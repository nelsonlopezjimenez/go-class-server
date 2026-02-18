// Package update implements update functions to ensure
// updated content on the local system.
package update

// Rocky Connor 420711

import (
	"fmt"
	"localhost/CIS/modules/command"
	"localhost/CIS/modules/external"

	"localhost/CIS/modules/util"
	"log"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
)

type NetworkPinger struct {
	Url     string
	Timeout int64
}

type IPInfo struct {
	IP          net.IP
	isConnected bool
}

type GitError struct {
	Message      string
	FailedCmd    []string
	ErrorWrapped error
}

var websitesPath = util.GetOSPaths()["websites"]

// Creates a logger instance specifically for the update functions to inform user of update related events
var updateLogger = log.New(os.Stdout, "[Updater] ", log.Ltime)
var ipData IPInfo

// Update
//
// Method calls checkForDependencies and initialize a ticker that fires the
// git pull command every np.Timeout minutes to update the monorepo and prints
// the results of the command to the console. Main.go calls this method as a go routine.
func (np NetworkPinger) Update() {
	ipData = GetLocalIP()
	if ipData.isConnected {
		updateLogger.Println("IP Address is: ", ipData.IP.To4().String())
		err := updateClassResources()
		if err != nil {
			updateLogger.Printf("Class Content failed to update: %v", err)
		}
	}

	interval := time.Minute
	if gin.Mode() == "debug" {
		interval = time.Second
	}
	checkInterval := time.NewTicker(time.Duration(np.Timeout) * interval)
	hasCheckedDeps := false
	if ipData.isConnected {
		err := checkForDependencies(np.Url)
		if err != nil {
			updateLogger.Println("[WARN]: Checking for dependencies failed:", err)
		} else {
			hasCheckedDeps = true
		}
		err = updateClassResources()
		if err != nil {
			updateLogger.Printf("Class Content failed to update: %v", err)
		}
	}
	for range checkInterval.C {
		ipData = GetLocalIP()

		if ipData.isConnected {
			fmt.Println(time.Now())
			if hasCheckedDeps {
				// external.UpdateSiteMetaData()
				err := updateClassResources()
				if err != nil {
					updateLogger.Printf("Class Content failed to update: %v", err)

				}
			}
		}

	}

}

// checkForDependencies
//
// Checks if directory named 'data' exists. If it does not, it runs the git clone
// command to clone the monorepo from Gitea.
func checkForDependencies(url string) error {
	checkUserConfig()
	_, dirErr := os.Stat(util.GetDepPath())
	if dirErr != nil {
		if os.IsNotExist(dirErr) {
			output, err := external.GitClone(".", url+"/ClassroomResources/ClassServerResources.git", "ClassServerResources")
			if err != nil {
				updateLogger.Println(err)
			}
			updateLogger.Printf("%s", output)
		} else {
			return fmt.Errorf("could not create ClassServerResources: %w", dirErr)
		}

	}
	_, websitesErr := os.Stat(websitesPath)
	if websitesErr != nil {
		fmt.Println("No websites folder")
		if os.IsNotExist(websitesErr) {
			err := os.Mkdir("C:/websites", 0755)
			if err != nil {
				updateLogger.Println("Creating websites dir:", err)
			}
		} else {
			return fmt.Errorf("could not create websites directory: %w", websitesErr)
		}
	}
	return nil
}

func updateClassResources() error {
	out, err := external.GitPull(util.GetDepPath())
	if err != nil {
		return err
	}
	updateLogger.Printf("%s", out)

	return nil
}

// GetLocalIP
//
// Checks for network interface other than localhost and returns a struct
// with the ip address and bool value. This allows for checking for network
// connection without sending get requests to the Gitea server over and over
func GetLocalIP() IPInfo {

	Info := IPInfo{nil, false}
	networkInterfaces, err := net.Interfaces()
	if err != nil {
		updateLogger.Println("[WARN]: Could not get network interfaces")
	}

	for _, networkInterface := range networkInterfaces {
		addrs, err := networkInterface.Addrs()
		if err != nil {
			updateLogger.Println("[WARN]: Could not get IP addresses")
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					Info = IPInfo{ipnet.IP, true}
				}
			}
		}
	}
	return Info
}

func (ge GitError) Error() string {
	retryCommand(ge)
	return fmt.Sprintf("%s. Are you connected to the dock?", ge.Message)
}

func retryCommand(ge GitError) {
	updateLogger.Println("retrying...")
	retryCmd := exec.Command(ge.FailedCmd[0], ge.FailedCmd[1:]...)

	for i := range 5 {
		time.Sleep(time.Duration(10^i*2) * time.Millisecond)

		out, err := retryCmd.CombinedOutput()
		if err != nil {
			updateLogger.Printf("Retry %v failed.", i+1)
		} else {
			updateLogger.Printf("%s, err: %v", out, err)
			break
		}
	}
}

func checkUserConfig() {
	defer func() {
		if err := recover(); err != nil {
			updateLogger.Fatal("A serious issue has occurred:", err)
		}
	}()
	_, err := exec.LookPath("git")
	if err != nil {
		panic("Git is not installed on your system.")
	}
	userName := os.Getenv("USERNAME")
	for Key, Value := range map[string]string{
		"user.name":  userName,
		"user.email": fmt.Sprintf("%s@edcc.edu", userName),
	} {
		err := command.CheckConfigKeyIsSet(Key, Value)
		if err != nil {
			panic("Failed to set git username and email")
		}
	}
}
