// services
//
// This package is for the handling of microservices for resources
// such as the video viewer backend.
package services

import (
	"fmt"
	"localhost/CIS/modules/logger"
	"localhost/CIS/modules/util"
	"os"
	"os/exec"
	"sync"
)

// StartVideoService
//
// Handles the video microservice and ensures shutdown of the microservice
// on server shutdown.
func StartVideoService(stopService chan struct{}, wg *sync.WaitGroup) {
	// Ensure the goroutine is removed from the wait group on return
	defer wg.Done()
	// logger.Log(logger.WarnLevel, "Not implemented yet...")
	//Create the exec.Cmd with arguments
	cmd := exec.Command("node", "server.js")
	// Load the path stored in env
	cmd.Dir = util.LoadEnv("VIDEO_VIEWER_PATH")
	// Env for the new process being started
	// TODO: Do this better
	cmd.Env = []string{"MONGO_URI=mongodb://localhost:27017/videosMasterDb", "PORT=3210"}
	// Start the process
	logger.Log(logger.InfoLevel, "Starting the video service")
	err := cmd.Start()
	if err != nil {
		if err := cmd.Process.Signal(os.Kill); err != nil {
			logger.Log(logger.ErrorLevel, fmt.Sprintf("Err sending kill sig: %v", err))
		}
		logger.Log(logger.ErrorLevel, fmt.Sprintf("Error from cmd is: %v", err))
	}
	// logger.Log(logger.InfoLevel, string(output))

	for {

		select {
		case <-stopService:
			// Case listens on stopService channel and executes the
			// following code on recept
			// TODO: Ensure the Express server gracefully shuts down
			// <-stopService
			logger.Log(logger.InfoLevel, "Stopping the video service")
			return
		}
	}
}
