// Rocky Connor 420711
package CIS

import (
	"fmt"
	command "localhost/CIS/modules/cmd"
	"localhost/CIS/modules/external"
	"localhost/CIS/modules/util"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type NetworkPinger struct {
	Url     string
	Timeout int64
}

type IPInfo struct {
	IP          string
	isConnected bool
}

type GitError struct {
	Message      string
	FailedCmd    []string
	ErrorWrapped error
}

var websitesPath = util.GetOSPaths().Websites

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
		err := external.UpdateSiteMetaData()
		if err != nil {
			updateLogger.Println(err)
		}
		err = updateClassResources()
		if err != nil {
			fmt.Println("I'm an error!")
			updateLogger.Println(err)
		}
	}

	interval := time.Minute
	if gin.Mode() == "debug" {
		interval = time.Second
	}
	checkInterval := time.NewTicker(time.Duration(np.Timeout) * interval)
	hasCheckedDeps := false
	if ipData.isConnected {
		err := external.UpdateSiteMetaData()
		if err != nil {
			updateLogger.Println(err)
		}
		checkForDependencies(np.Url)
		err = updateClassResources()
		if err != nil {
			updateLogger.Println(err)
		}
		hasCheckedDeps = true
	}
	for range checkInterval.C {
		ipData = GetLocalIP()

		if ipData.isConnected {
			fmt.Println(time.Now())
			if hasCheckedDeps {
				external.UpdateSiteMetaData()
				err := updateClassResources()
				if err != nil {
					updateLogger.Println(err)
				}
			}
		}

	}

}

// checkForDependencies
//
// Checks if directory named 'data' exists. If it does not, it runs the git clone
// command to clone the monorepo from Gitea.
func checkForDependencies(url string) {
	checkUserConfig()
	cwd, err := os.Getwd()
	if err != nil {
		updateLogger.Panicln("Cannot get CWD:", err)
	}
	cwd = strings.Replace(cwd, "\\", "/", -1)
	_, dirErr := os.Stat(cwd + "/data")
	if dirErr != nil {
		if os.IsNotExist(dirErr) {
			output, err := GitClone(".", url+"/ClassroomResources/ClassServerResources.git", "data")
			if err != nil {
				updateLogger.Println(err)
			}
			updateLogger.Printf("%s", output)
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
		}
	}
}

func updateClassResources() error {
	out, err := GitPull("./data")
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

	Info := IPInfo{"", false}
	networkInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("[IP Finder]: Could not get interfaces")
	}

	for _, networkInterface := range networkInterfaces {
		addrs, err := networkInterface.Addrs()
		if err != nil {
			fmt.Println("[IP Finder]: Could not get addresses")
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					fmt.Println()
					Info = IPInfo{ipnet.IP.String(), true}
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
	userName := os.Getenv("USERNAME")
	for Key, Value := range map[string]string{
		"user.name":  userName,
		"user.email": fmt.Sprintf("%s@edcc.edu", userName),
	} {
		err := command.CheckConfigKeyIsSet(Key, Value)
		if err != nil {
			fmt.Println("[cUC]:", err)
		}
	}
}
