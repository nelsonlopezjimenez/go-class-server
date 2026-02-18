package main

// Rocky Connor 420711
// The package the code belongs to must be declared at the top of the file
// If the package is to be the entry point, it must be declared as main

// All external packages must be declared next
// Every package that is declared in the imports must be used
import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"localhost/CIS/modules/controller"
	"localhost/CIS/modules/update"
	"localhost/CIS/modules/util"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var filePath = util.GetOSPaths()

// The main package must contain a main function which will be executed on run
func main() {
	godotenv.Load() //? load .env file if it exists

	//! All variables that could be modified by .env should use
	//! util.LoadEnv(key string) and not try to access the environment directly
	releaseVersion := util.LoadEnv("RELEASE_VERSION")
	releaseDate := util.LoadEnv("RELEASE_DATE")
	PORT := util.LoadEnv("CIS_CLASS_SERVER_PORT")

	// Creates a new logger instance to display info from server
	serverLog := log.New(os.Stdout, "[Server] ", log.LstdFlags)
	serverLog.Println("CIS Class Server", releaseVersion, "released on", releaseDate)

	// The router var creates the default gin engine instance
	if !util.IsDevelopment() {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	// server var changed to create an instance of http.Server and
	// the default gin instance is passed as a handler to take advantage
	// of the http.Server.Shutdown method to facilitate graceful shutdown
	server := &http.Server{
		Addr:    ":" + PORT,
		Handler: router,
	}
	//
	router.Use(controller.IsAuthorized)
	// router.Use(controller.CheckForSiteExt)
	api := router.Group("/api")
	websitesGroup := router.Group("/websites")
	websitesGroup.Use(controller.CheckForSiteExt)
	// The following line defines a static asset folder
	fsys, err := util.GetFileSystemHandler()
	if err != nil {
		serverLog.Println("there was an error in the embedded fs:", err)
	}
	router.StaticFS("/assets", fsys)
	router.Static("/images", util.GetDepPath()+"/images")
	_, staticErr := os.Stat(filePath["static"])
	if staticErr != nil {

		if os.IsNotExist(staticErr) {
			fmt.Println("staticErr:", staticErr)
			err := os.Mkdir(filePath["static"], 0777)
			if err != nil {
				serverLog.Println("Error creating /static:", err)
			}
		}
	}

	util.CheckForSymlink("Public")
	util.CheckForSymlink("CANVAS_FILE_CACHES")
	util.CheckForSymlink("Videos")

	router.StaticFS("/static", gin.Dir(filePath["serverPath"]+"/static", true))
	websitesGroup.StaticFS("/", gin.Dir(filePath["websites"], true))

	Gitea := update.NetworkPinger{Url: util.LoadEnv("UPDATE_IP"), Timeout: 10}

	// Goroutine to check for lesson repo and updates if there is a connection
	go Gitea.Update()
	// This middleware function returns the requested offline website to the client
	router.GET("/", controller.SendIndex)
	// router.GET("/websites/*url", controller.GetWebsite)
	router.POST("/websites/try.w3schools.com/*path", controller.HandleUserCode)
	// router.GET("/websites", controller.CheckForSiteExt)
	router.GET("/:allOther/*any", controller.SendIndex)
	router.GET("/raw/lessons/:mdFile/:lesson", controller.SendOneLesson)

	api.GET("/information/:infoPage", controller.SendInformationPage)
	api.GET("/lessons", controller.SendLessonsList)
	api.GET("/links", controller.SendLinks)
	api.GET("/git/:command/:submodule", controller.HandleWebsiteManagement)

	if !util.IsDevelopment() {
		url := "http://localhost:" + PORT
		util.OpenBrowser(url)
	}

	// Start the server in a goroutine as to not block the graceful
	// shutdown logic after the call to server.ListenAndServe
	go func() {
		// Starts the server on the specified port
		serverLog.Printf("Server running on %v", PORT)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverLog.Println("Error in the server:", err)
		}
	}()

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
