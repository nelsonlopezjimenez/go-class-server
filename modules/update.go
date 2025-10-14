// Rocky Connor 420711
package CIS

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"regexp"
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

var websitesPath = GetOSPaths().Websites

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
		err := updateClassResources()
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
		checkForDependencies(np.Url)
		err := updateClassResources()
		if err != nil {
			updateLogger.Println(err)
		}
		err = pullWebsitesSuperproject()
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
				err := updateClassResources()
				if err != nil {
					updateLogger.Println(err)
				}
				if gin.Mode() == "release" {
					updateLogger.Println("Running pullWebsitesSuperproject()")
					siteErr := pullWebsitesSuperproject()
					if siteErr != nil {
						updateLogger.Println(siteErr)
					}
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
	cwd, err := os.Getwd()
	if err != nil {
		updateLogger.Panicln("Cannot get CWD:", err)
	}
	cwd = strings.Replace(cwd, "\\", "/", -1)
	_, dirErr := os.Stat(cwd + "/data")
	if dirErr != nil {
		if os.IsNotExist(dirErr) {
			Gitea := exec.Command("git", "clone", "-b", "main", url+"/ClassroomResources/ClassServerResources.git", "./data")
			updateLogger.Println("data directory does not Exist.")
			updateLogger.Println("Creating it now.")
			updateLogger.Println("Attempting to clone data")
			_, cloneErr := Gitea.CombinedOutput()
			if cloneErr != nil {
				updateLogger.Println(cloneErr)
			}
		}

	}
	_, websitesErr := os.Stat(websitesPath)
	if websitesErr != nil {
		fmt.Println("No websites folder")
		if os.IsNotExist(websitesErr) {
			superErr := getWebsitesSuperproject()
			if superErr != nil {
				updateLogger.Println(superErr)

			}

		}
	}
}

// getWebsitesSuperproject
//
// Issues git clone command to clone websites super project.
func getWebsitesSuperproject() error {
	websitesSuper := exec.Command("git", "clone", "http://192.168.1.47:3000/OfflineWebsites/websites.git")
	websitesSuper.Dir = GetRoot()
	out, err := websitesSuper.CombinedOutput()
	if err != nil {
		// updateLogger.Println("Error cloning websites Superproject!!:", err)
		return GitError{"Could not clone the websites superproject. Are you on the dock?", websitesSuper.Args, err}
	}

	updateLogger.Println("Cloned websites superproject:", string(out))
	return nil

}

// GetWebsiteModule
//
// Executes git submodule command to initiate clone of selected submodule
func GetWebsiteModule(url string) (string, error) {
	updateCmd := exec.Command("git", "submodule", "update", "--init", "--remote", url)
	updateCmd.Dir = websitesPath
	out, err := updateCmd.CombinedOutput()
	if err != nil {

		return "", GitError{"Failed to get " + url + ". Are you connected to the dock?", updateCmd.Args, err}
	}

	return string(out), nil
}

// pullWebsitesSuperproject
//
// Issues the git pull command to update the websites superproject
// and then calls submoduleUpdateAll()
func pullWebsitesSuperproject() error {
	websitesPull := exec.Command("git", "pull", "--force", "http://192.168.1.47:3000/OfflineWebsites/websites.git")
	websitesPull.Dir = websitesPath
	out, err := websitesPull.CombinedOutput()
	if err != nil {
		return GitError{"Could not update the superproject", websitesPull.Args, err}
	}

	updateLogger.Println(string(out))
	submoduleUpdateAll()

	return nil

}

// updateClassResources
//
// executes git pull command to get updates to the class resources repo
func updateClassResources() error {
	gitPull := exec.Command("git", "pull", "--force", "origin", "main")
	gitPull.Dir = "./data"
	updateLogger.Println("Checking for class content")
	err := gitPull.Run()
	if err != nil {
		updateLogger.Println("err:", err)
		return GitError{"Failed to update lessons.", gitPull.Args, err}

	}
	return nil
}

// updateSubmodule
//
// updateSubmodule should run for every submodule in the websites folder
// executes "git submodule update --remote" on supplied repository.
func updateSubmodule(path string) error {
	subUpdate := exec.Command("git", "submodule", "update", "--remote", path)
	subUpdate.Dir = websitesPath
	output, err := subUpdate.CombinedOutput()
	if err != nil {
		updateLogger.Println("Error updating website", path+":", err)
		return GitError{
			"Failed to update " + path,
			subUpdate.Args,
			err,
		}
	}

	updateLogger.Println(path+":", string(output))
	return nil
}

// submoduleUpdateAll
//
// Runs git submodule status on the websites superproject and
// finds all statuses beginning with "+". Then this string is
// split at the spaces to isolate the submodule name that can
// then be passed to updateSubmodule()
func submoduleUpdateAll() {

	status := exec.Command("git", "submodule", "status")
	status.Dir = "C:/websites"

	out, _ := status.Output()

	findNeedsUpdate, _ := regexp.Compile(`\+.+`)
	foundNeedsUpdate := findNeedsUpdate.FindAll(out, -1)
	if foundNeedsUpdate != nil {
		for _, line := range foundNeedsUpdate {
			submodule := strings.Split(string(line), " ")

			updateSubmodule(submodule[1])
		}
	} else {
		updateLogger.Println("All downloaded websites are up to date!")
	}
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
