// Package controller implements route handlers to handle calls to the HTTP
// server for the various routes.
package controller

// Rocky Connor 420711

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"localhost/CIS/modules/external"
	"localhost/CIS/modules/logger"
	"localhost/CIS/modules/util"

	"github.com/gin-gonic/gin"
)

var filePath = util.GetOSPaths()

func SendIndex(ctx *gin.Context) {
	// route function for handling requests to the root
	index, err := util.GetIndex()
	if err != nil {
		fmt.Println(err)
		// serverLog.Panicln(err)
	}
	fmt.Fprintf(ctx.Writer, "%s", index)
}

func SendOneLesson(ctx *gin.Context) {
	subdir := ctx.Param("mdFile")
	lessonName := ctx.Param("lesson")
	// creates a fs.FS  for the information directory
	fsys := os.DirFS(util.GetDepPath() + "/markdown/lessons/" + subdir)
	// Opens the requested markdown file
	file, err := fs.ReadFile(fsys, lessonName+".md")
	if err != nil {
		logger.Log(logger.WarnLevel, fmt.Sprintf("could not get %v: %v", lessonName, err))
		// serverLog.Panicln("There was an error getting the requested file:", err)
		ctx.JSON(500, fmt.Errorf("could not get %v: %w", lessonName, err))
		return
	}

	ctx.JSON(200, string(file))
}

func SendInformationPage(ctx *gin.Context) {
	reqInfoPage := ctx.Param("infoPage")
	// creates a fs.FS  for the information directory
	fsys := os.DirFS(util.GetDepPath() + "/markdown/information")
	// removes the leading / from the wildcard param
	reqInfoPage = strings.Replace(reqInfoPage, "/", "", 1)
	// Opens the requested markdown file
	file, err := fs.ReadFile(fsys, reqInfoPage+".md")
	if err != nil {
		logger.Log(logger.WarnLevel, fmt.Sprintf("There was an error getting the requested file: %v", err))
	}

	ctx.JSON(200, string(file))
}

func SendLessonsList(ctx *gin.Context) {
	testList := util.MakeRootDir(util.GetDepPath() + "/markdown/lessons")

	testSlice := map[string][]string{}

	testList.RecursiveSearchByExt(".md", func(path string, fileName string) {

		testSlice[path] = append(testSlice[path], fileName)
	})
	ctx.JSON(200, testSlice)
}

func SendLinks(ctx *gin.Context) {
	websiteInfoSlice, err := external.SendAllSites()
	if err != nil {
		ctx.String(500, err.Error())
		return
	}

	// Sends response  json data to the client
	ctx.JSON(200, websiteInfoSlice)

}

// This allows the available w3schools examples to execute
// A browser extension is also required to route the post
// request from W3S to the localhost
func HandleUserCode(ctx *gin.Context) {
	code := ctx.Request.FormValue("code")
	// code2 := ctx.Request.FormValue("code2")
	// code3 := ctx.Request.FormValue("code3")
	// codeInput := ctx.Request.FormValue("codeInput")
	lang := ctx.Request.FormValue("lang")
	// RunCode is a function imported from the CIS module
	output, err := util.RunUserCode(code, lang)
	if err != nil {
		logger.Log(logger.DebugLevel, fmt.Sprintf("Error running user supplied code: %v", err))
	}

	fmt.Fprintf(ctx.Writer, "%s", output)

}

func HandleWebsiteManagement(ctx *gin.Context) {
	type ReturnOutput map[string]string
	submodule := ctx.Param("submodule")
	command := ctx.Param("command")
	var consoleOutput string
	switch command {
	case "update":
		consoleOutputBytes, err := external.GitPull(filePath["websites"] + "/" + submodule)
		if err != nil {
			logger.Log(logger.WarnLevel, fmt.Sprintf("Could not download new files: %v", err))

			ctx.JSON(500, err.Error())
			return
		}
		err = external.UpdateSingleInfo(submodule, false)
		if err != nil {
			logger.Log(logger.WarnLevel, fmt.Sprintf("Could not update site data: %v", err.Error()))

			ctx.JSON(500, err.Error())
			return
		}
		consoleOutput = string(consoleOutputBytes)
	case "install":
		consoleOutputBytes, err := external.GitClone(filePath["websites"], util.LoadEnv("GIT_INSTALL_ADDR")+submodule+".git", submodule)
		if err != nil {
			logger.Log(logger.WarnLevel, fmt.Sprintf("Could not download %v: %v", submodule, err))

			ctx.JSON(500, err.Error())
			return
		}

		err = external.UpdateSingleInfo(submodule, false)
		if err != nil {
			logger.Log(logger.WarnLevel, fmt.Sprintf("Could not update %v: %v\n", submodule, err.Error()))
			ctx.JSON(500, err.Error())
			return
		}

		consoleOutput = string(consoleOutputBytes)
	case "delete":
		// delete specified domain dir
		err := util.DeleteSite(filePath["websites"] + "/" + submodule)
		if err != nil {
			consoleOutput = err.Error()
			logger.Log(logger.WarnLevel, fmt.Sprintf("Could not delete %v: %v ", submodule, err))
		}
		err = external.UpdateSingleInfo(submodule, true)
		if err != nil {
			logger.Log(logger.WarnLevel, fmt.Sprintf("Could not update state of %v: %v", submodule, err.Error()))
			ctx.JSON(500, err.Error())
			return
		}
	default:
		ctx.Status(403)
		return
	}

	ctx.JSON(200, ReturnOutput{"output": string(consoleOutput)})
}
