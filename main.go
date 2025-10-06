// Rocky Connor 420711
// The package the code belongs to must be declared at the top of the file
// If the package is to be the entry point, it must be declared as main
package main

// All external packages must be declared next
// Every package that is declared in the imports must be used
import (
	// "bytes"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"strings"

	CIS "localhost/CIS/modules"

	"github.com/gin-gonic/gin"
)

// Defines the port the server binds to
var (
	port = flag.String("p", "22022", "Sets the port for the server. Default is 22022.")
	dev  = flag.Bool("dev", false, "Runs the server in dev mode. Displays debug messages. Default is false.")
)

var (
	releaseVersion = "v1.1.2"
	releaseDate = "10/2/2025"
)

var updateIP string

func usage() {
	fmt.Println("usage: ClassServer -flag options")
}

// The main package must contain a main function which will be executed on run
func main() {
	// Here we parse the flags in case of user defined options
	flag.Usage = usage
	flag.Parse()
	isDev := *dev

	// Creates a new logger instance to display info from server
	serverLog := log.New(os.Stdout, "[Server] ", log.LstdFlags)
	serverLog.Println("CIS Class Server", releaseVersion, "released on", releaseDate )
	if isDev {
		serverLog.Println("Running in dev mode.")
		updateIP = "http://localhost:3000"
		serverLog.Println("The update  URL is:", updateIP)
	}

	if !isDev {
		gin.SetMode(gin.ReleaseMode)
		updateIP = "http://192.168.1.28:3000"
		serverLog.Println("The update  URL is:", updateIP)
	}
	// The server var creates the default gin engine instance
	server := gin.Default()
	api := server.Group("/api")
	// The following line defines a static asset folder
	fsys, err := CIS.GetFileSystemHandler()
	if err != nil {
		serverLog.Println("there was an error in the embedded fs:", err)
	}
	server.StaticFS("/assets", fsys)
	server.Static("/images", "./data/images")
	_, staticErr := os.Stat("C:/Users/Public/ClassServer/static")
	if staticErr != nil {
		if os.IsNotExist(staticErr) {
			err := os.Mkdir("C:/Users/Public/classServer/static/", 0777)
			if err != nil {
				serverLog.Println("Error creating /static:", err)
			} else {
				cfcErr := os.Symlink("C:/Users/Public/CANVAS_FILE_CACHES/", "C:/Users/Public/classServer/static/CANVAS_FILE_CACHES")
				if cfcErr != nil {
					fmt.Println("cfc:", cfcErr)
				}
				videoErr := os.Symlink("C:/Users/Public/Videos/", "C:/Users/Public/classServer/static/Videos")
				if videoErr != nil {
					fmt.Println("video:", videoErr)
				}
			}
		}
	}

	server.Static("/static", "C:/Users/Public/ClassServer/static")

	Gitea := CIS.NetworkPinger{Url: updateIP, Timeout: 10}
	// Goroutine to check for lesson repo and updates if there is a connection
	go Gitea.Update()
	// This middleware function returns the requested offline website to the client
	server.GET("/websites/*url", func(ctx *gin.Context) {
		// ctc.Param returns the wildcard value in the url path
		param := ctx.Param("url")
		ctx.File("C:/websites" + param)
	})

	// This allows the available w3schools examples to execute
	// A browser extension is also required to route the post
	// request from W3S to the localhost
	server.POST("/websites/try.w3schools.com/*path", func(ctx *gin.Context) {
		code := ctx.Request.FormValue("code")
		// code2 := ctx.Request.FormValue("code2")
		// code3 := ctx.Request.FormValue("code3")
		// codeInput := ctx.Request.FormValue("codeInput")
		lang := ctx.Request.FormValue("lang")
		// RunCode is a function imported from the CIS module
		output, err := CIS.RunCode(code, lang)
		if err != nil {
			serverLog.Println("there was an error: ", err)
		}

		fmt.Fprintf(ctx.Writer, "%s", output)

	})

	server.GET("/", func(ctx *gin.Context) {
		// route function for handling requests to the root
		index, err := CIS.GetIndex()
		if err != nil {
			serverLog.Panicln(err)
		}
		// serverLog.Println("redirect header is:", ctx.Request.Header)
		if len(ctx.Request.Header.Get("Referer")) > 0 {
			ctx.Header("redirect", ctx.Request.Header.Get("Referer"))
		}
		fmt.Fprintf(ctx.Writer, "%s", index)
	})

	server.GET("/:allOther", func(ctx *gin.Context) {
		// ctx.Redirect(301, "/")
		index, err := CIS.GetIndex()
		if err != nil {
			serverLog.Panicln(err)
		}
		fmt.Fprintf(ctx.Writer, "%s", index)

	})

	server.GET("/lessons/:mdFile", func(ctx *gin.Context) {
		// renders markdown files into HTML and sends to client

		lessonPage := ctx.Param("mdFile")

		// if no markdown file is requested, returns information.md by default
		if lessonPage == "/" {
			lessonPage = "index"
		}

		// creates a fs.FS  for the markdown directory
		fsys := os.DirFS("./data/markdown")
		// removes the leading / from the wildcard param
		lessonPage = strings.Replace(lessonPage, "/", "", 1)
		// Opens the requested markdown file
		file, err := fs.ReadFile(fsys, lessonPage+".md")
		if err != nil {
			serverLog.Panicln("There was an error getting the mardown file:", err)
		}

		// Parses the requested file from markdown to HTML
		// Sends the parsed HTML to the client as JSON data
		ctx.JSON(200, string(file))
	})

	server.GET("/lessons/:mdFile/:lesson", func(ctx *gin.Context) {
		subdir := ctx.Param("mdFile")
		lessonName := ctx.Param("lesson")
		// creates a fs.FS  for the information directory
		fsys := os.DirFS("./data/markdown/lessons/" + subdir)
		// Opens the requested markdown file
		file, err := fs.ReadFile(fsys, lessonName+".md")
		if err != nil {
			serverLog.Panicln("There was an error getting the requested file:", err)
		}

		ctx.JSON(200, string(file))
	})

	server.GET("/raw/lessons/:mdFile/:lesson", func(ctx *gin.Context) {
		subdir := ctx.Param("mdFile")
		lessonName := ctx.Param("lesson")
		// creates a fs.FS  for the information directory
		fsys := os.DirFS("./data/markdown/lessons/" + subdir)
		// Opens the requested markdown file
		file, err := fs.ReadFile(fsys, lessonName+".md")
		if err != nil {
			serverLog.Panicln("There was an error getting the requested file:", err)
		}

		ctx.JSON(200, string(file))
	})

	api.GET("/information/:infoPage", func(ctx *gin.Context) {
		reqInfoPage := ctx.Param("infoPage")
		// creates a fs.FS  for the information directory
		fsys := os.DirFS("./data/markdown/information")
		// removes the leading / from the wildcard param
		reqInfoPage = strings.Replace(reqInfoPage, "/", "", 1)
		// Opens the requested markdown file
		file, err := fs.ReadFile(fsys, reqInfoPage+".md")
		if err != nil {
			serverLog.Panicln("There was an error getting the requested file:", err)
		}

		ctx.JSON(200, string(file))
	})

	api.GET("/lessons", func(ctx *gin.Context) {
		testList := CIS.RootDir{Root: "./data/markdown/lessons"}

		testSlice := map[string][]string{}

		// fmt.Println(testSlice)
		testList.RecursiveSearch(".md", func(path string, fileName string) {

			testSlice[path] = append(testSlice[path], fileName)
		})
		// fmt.Println(testSlice)
		ctx.JSON(200, testSlice)
	})

	api.GET("/links", func(ctx *gin.Context) {

		// Remember: the property names must be UPPERCASE in order to be exported
		// Type Outbound defines what the structure should contain.

		newList := CIS.RootDir{Root: "C:/websites"}
		newListTest := newList.ListBuilder()
		type testStruct struct {
			Domain    string
			IndexPath string
			Installed bool
			IsCurrent bool
		}
		structSlice := []testStruct{}
		for _, item := range newListTest {

			isHidden := strings.HasPrefix(item, ".")
			if !isHidden {
				var index string
				currentDir := CIS.RootDir{Root: "C:/websites/" + item}
				currentDir.FindIndex(func(path, fileName string) {
					index = path + "/" + fileName
				})
				pageUpdated := true
				isInstalled := false
				if index != "" {
					isInstalled = true
				}
				indexStruct := testStruct{item, index, isInstalled, pageUpdated}
				structSlice = append(structSlice, indexStruct)
			}
		}

		// Sends response  json data to the client
		ctx.JSON(200, structSlice)

	})

	api.GET("/git/update/:submodule", func(ctx *gin.Context) {
		submodule := ctx.Param("submodule")

		gitOut := CIS.GetWebsiteModule(submodule)
		ctx.JSON(200, gitOut)
	})

	userHome, err := os.UserHomeDir()
	if err != nil {
		serverLog.Println("there was an error getting the user's home dir:", err)
	}

	userHome = strings.ReplaceAll(userHome, "\\", "/")

	if !isDev {
		cmd := exec.Command(userHome+"/AppData/Local/Google/Chrome/Application/chrome.exe", "http://localhost:"+*port)
		if err := cmd.Start(); err != nil {
			serverLog.Println("Error opening Chrome:", err)
		}
	}

	// Starts the server on the specified port
	server.Run(":" + *port)

}
