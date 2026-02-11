// Rocky Connor 420711
// The package the code belongs to must be declared at the top of the file
// If the package is to be the entry point, it must be declared as main
package main

// All external packages must be declared next
// Every package that is declared in the imports must be used
import (
	// "bytes"

	"context"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	CIS "localhost/CIS/modules"
	"localhost/CIS/modules/external"
	"localhost/CIS/modules/util"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Defines the port the server binds to
var (
	port = flag.String("p", "22022", "Sets the port for the server. Default is 22022.")
	dev  = flag.Bool("dev", false, "Runs the server in dev mode. Displays debug messages. Default is false.")
)

var updateIP string

var filePath = util.GetOSPaths()

func usage() {
	fmt.Println("usage: ClassServer -flag options")
}

// The main package must contain a main function which will be executed on run
func main() {
	godotenv.Load() //? load .env file if it exists

	//! All variables that could be modified by .env should use
	//! util.LoadEnv(key string) and not try to access the environment directly
	releaseVersion := util.LoadEnv("RELEASE_VERSION")
	releaseDate := util.LoadEnv("RELEASE_DATE")

	// Here we parse the flags in case of user defined options
	flag.Usage = usage
	flag.Parse()
	isDev := *dev

	// Creates a new logger instance to display info from server
	serverLog := log.New(os.Stdout, "[Server] ", log.LstdFlags)
	serverLog.Println("CIS Class Server", releaseVersion, "released on", releaseDate)
	if isDev {
		serverLog.Println("Running in dev mode.")
		updateIP = util.LoadEnv("UPDATE_IP")
		serverLog.Println("The update  URL is:", updateIP)
	}

	if !isDev {
		gin.SetMode(gin.ReleaseMode)
		updateIP = util.LoadEnv("UPDATE_IP")
		serverLog.Println("The update  URL is:", updateIP)
	}

	// The router var creates the default gin engine instance
	// Todo: Move all gin related logic to its own module
	router := gin.Default()
	// server var changed to create an instance of http.Server and
	// the default gin instance is passed as a handler to take advantage
	// of the http.Server.Shutdown method to facilitate graceful shutdown
	server := &http.Server{
		Addr:    ":" + *port,
		Handler: router,
	}

	// This middleware function prevents the access of the students' server
	// to anyone on the outside network. This prevents the possibility of data
	// transfer from system to system as prohibited by DOC.
	// Note: This can only be controlled via the source code. After compile,
	// This cannot be circumvented from the executable itself.
	router.Use(func(ctx *gin.Context) {
		switch ctx.RemoteIP() {
		case "::1", "127.0.0.1":
			ctx.Next()
		default:
			ctx.String(403, "I'm sorry. Your are not authorized to see this.")
			ctx.Abort()
		}
	})
	api := router.Group("/api")
	// The following line defines a static asset folder
	fsys, err := util.GetFileSystemHandler()
	if err != nil {
		serverLog.Println("there was an error in the embedded fs:", err)
	}
	router.StaticFS("/assets", fsys)
	router.Static("/images", util.GetDepPath()+"/images")
	_, staticErr := os.Stat(filePath.ServerPath + "/static")
	if staticErr != nil {

		if os.IsNotExist(staticErr) {
			fmt.Println("staticErr:", staticErr)
			err := os.Mkdir(filePath.ServerPath+"/static", 0777)
			if err != nil {
				serverLog.Println("Error creating /static:", err)
			}
		}
	}
	_, publicErr := os.Stat(filePath.ServerPath + "/static/Public")
	if os.IsNotExist(publicErr) {
		publicErr := os.Symlink(filePath.Public, filePath.ServerPath+"/static/Public")
		if publicErr != nil {
			fmt.Println("Public:", publicErr)
		}
	}
	_, cfcErr := os.Stat(filePath.ServerPath + "/static/CANVAS_FILE_CACHES")
	if os.IsNotExist(cfcErr) {
		cfcErr := os.Symlink(filePath.CFC, filePath.ServerPath+"/static/CANVAS_FILE_CACHES")
		if cfcErr != nil {
			fmt.Println("cfc:", cfcErr)
		}
	}
	_, videoErr := os.Stat(filePath.ServerPath + "/static/Videos")
	if os.IsNotExist(videoErr) {
		videoErr := os.Symlink(filePath.Videos, filePath.ServerPath+"/static/Videos")
		if videoErr != nil {
			fmt.Println("video:", videoErr)
		}
	}

	router.StaticFS("/static", gin.Dir(filePath.ServerPath+"/static", true))

	Gitea := CIS.NetworkPinger{Url: updateIP, Timeout: 10}
	// Goroutine to check for lesson repo and updates if there is a connection
	go Gitea.Update()
	// This middleware function returns the requested offline website to the client
	router.GET("/websites/*url", func(ctx *gin.Context) {
		// ctc.Param returns the wildcard value in the url path
		param := ctx.Param("url")

		if !strings.HasSuffix(param, ".html") && !strings.HasSuffix(param, "/") {
			if !util.CheckExt(param) {
				param = param + ".html"
			}
		}

		if strings.HasSuffix(param, ".asp") {
			param = strings.Replace(param, ".asp", ".html", 1)
		}

		fmt.Println(param)
		ctx.File(filePath.Websites + "/" + param)
	})

	// This allows the available w3schools examples to execute
	// A browser extension is also required to route the post
	// request from W3S to the localhost
	router.POST("/websites/try.w3schools.com/*path", func(ctx *gin.Context) {
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

	router.GET("/", func(ctx *gin.Context) {
		// route function for handling requests to the root
		index, err := util.GetIndex()
		if err != nil {
			serverLog.Panicln(err)
		}
		fmt.Fprintf(ctx.Writer, "%s", index)
	})

	router.GET("/:allOther/*any", func(ctx *gin.Context) {
		index, err := util.GetIndex()
		if err != nil {
			serverLog.Panicln(err)
		}
		fmt.Fprintf(ctx.Writer, "%s", index)

	})

	router.GET("/raw/lessons/:mdFile/:lesson", func(ctx *gin.Context) {
		subdir := ctx.Param("mdFile")
		lessonName := ctx.Param("lesson")
		// creates a fs.FS  for the information directory
		fsys := os.DirFS(util.GetDepPath() + "/markdown/lessons/" + subdir)
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
		fsys := os.DirFS(util.GetDepPath() + "/markdown/information")
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
		testList := util.RootDir{Root: util.GetDepPath() + "/markdown/lessons"}

		testSlice := map[string][]string{}

		testList.RecursiveSearch(".md", func(path string, fileName string) {

			testSlice[path] = append(testSlice[path], fileName)
		})
		ctx.JSON(200, testSlice)
	})

	api.GET("/links", func(ctx *gin.Context) {
		websiteInfoSlice, err := external.SendAllSites()
		if err != nil {
			ctx.String(500, err.Error())
			return
		}

		// Sends response  json data to the client
		ctx.JSON(200, websiteInfoSlice)

	})

	api.GET("/git/:command/:submodule", func(ctx *gin.Context) {
		type ReturnOutput map[string]string
		submodule := ctx.Param("submodule")
		command := ctx.Param("command")
		serverLog.Println(command)
		var consoleOutput string
		switch command {
		case "update":
			consoleOutputBytes, err := CIS.GitPull(filePath.Websites + "/" + submodule)
			if err != nil {
				fmt.Printf("err.Error() pull: %v\n", err.Error())
				ctx.JSON(500, err.Error())
				return
			}
			err = external.UpdateSingleInfo(submodule)
			if err != nil {
				fmt.Printf("err.Error() upSing: %v\n", err.Error())
				ctx.JSON(500, err.Error())
				return
			}
			consoleOutput = string(consoleOutputBytes)
		case "install":
			consoleOutputBytes, err := CIS.GitClone(filePath.Websites, util.LoadEnv("GIT_INSTALL_ADDR")+submodule+".git", submodule)
			if err != nil {
				fmt.Println("[main]", err)

				ctx.JSON(500, err.Error())
				return
			}

			err = external.UpdateSingleInfo(submodule)
			if err != nil {
				fmt.Printf("err.Error() upSing: %v\n", err.Error())
				ctx.JSON(500, err.Error())
				return
			}

			consoleOutput = string(consoleOutputBytes)
		case "delete":
			// delete specified domain dir
			err := util.DeleteSite(filePath.Websites + "/" + submodule)
			if err != nil {
				consoleOutput = err.Error()
			}
		default:
			ctx.Status(403)
			return
		}

		ctx.JSON(200, ReturnOutput{"output": string(consoleOutput)})
	})

	if !isDev {
		url := "http://localhost:" + *port
		util.OpenBrowser(url)
	}

	// Start the server in a goroutine as to not block the graceful
	// shutdown logic after the call to server.ListenAndServe
	go func() {
		// Starts the server on the specified port
		svrErr := server.ListenAndServe()
		if svrErr != nil && !errors.Is(svrErr, http.ErrServerClosed) {
			serverLog.Println("Error in the server:", svrErr)
		}
	}()

	// TODO: Separate this into its own package/file
	// Everything beyond this line is logic for handling graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	serverLog.Println("Gracefully shutting down the server.")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		serverLog.Fatal("Server forced to shutdown:", err)
	}
	serverLog.Println("Server exiting")
}
