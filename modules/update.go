// Rocky Connor 420711
package CIS

import (
	"fmt"
	"log"
	"net/http"
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

// Creates a logger instance specifically for the update functions to inform user of update related events
var updateLogger = log.New(os.Stdout, "[Updater] ", log.Ltime)

// checkForServer
//
// This method is an example of a non-exported function. Functions are only exported if their name begins with a capital.
// Attempts to send a HTTP request to the url string defined the NetworkPinger.Url.
// If it receives a response, it returns true, otherwise it returns false if an exception is raised.
func (np NetworkPinger) checkForServer() bool {
	_, err := http.Get(np.Url)
	if err != nil {
		return false
	} else {
		return true
	}
}

// Update
//
// Method calls checkForDependencies and initialize a ticker that fires the
// git pull command every np.Timeout minutes to update the monorepo and prints
// the results of the command to the console.
func (np NetworkPinger) Update() {
	interval := time.Minute
	if gin.Mode() == "debug" {
		interval = time.Second
	}
	checkInterval := time.NewTicker(time.Duration(np.Timeout) * interval)
	hasCheckedDeps := false
	if np.checkForServer() {
		checkForDependencies(np.Url)
		updateClassResources()
		hasCheckedDeps = true
	}
	for range checkInterval.C {
		isAvailable := np.checkForServer()
		// websitePull := exec.Command("git", "pull", "--force", "origin", "main")

		if isAvailable {
			if hasCheckedDeps {
				updateClassResources()
				if gin.Mode() == "release" {
					siteErr := pullWebsitesSuperproject()
					if siteErr != nil {
						updateLogger.Println("Error pulling sites:", siteErr)
					}
				}

			} else {
				checkForDependencies(np.Url)
				hasCheckedDeps = true
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
			cloneTicker := time.NewTicker(30 * time.Second)
			updateLogger.Println("data directory does not Exist.")
			updateLogger.Println("Creating it now.")
			updateLogger.Println("Attempting to clone data")
			_, cloneErr := Gitea.CombinedOutput()
			if cloneErr != nil {
				updateLogger.Println("The clone failed, Trying again")
				updateLogger.Println(cloneErr)
				for range cloneTicker.C {
					_, recloneErr := Gitea.CombinedOutput()
					if recloneErr != nil {
						updateLogger.Println("The clone failed, Trying again")
						updateLogger.Println(cloneErr)
					} else {

						cloneTicker.Stop()
					}
				}

			}
		}

	}
	_, websitesErr := os.Stat("C:/websites")
	websitesTicker := time.NewTicker(30 * time.Second)
	if websitesErr != nil {
		fmt.Println("No websites folder")
		if os.IsNotExist(websitesErr) {
			superErr := getWebsitesSuperproject()
			if superErr != nil {
				updateLogger.Println("Failed to clone websites super project. Trying again...")
				for range websitesTicker.C {
					gayFish := getWebsitesSuperproject()
					if gayFish != nil {
						updateLogger.Println("Failed to clone websites super project. Trying again...")
					} else {
						websitesTicker.Stop()
					}
				}
			}

		}
	}
}

func getWebsitesSuperproject() error {
	websitesSuper := exec.Command("git", "clone", "http://192.168.1.47:3000/OfflineWebsites/websites.git")
	websitesSubmodules := exec.Command("git", "submodule", "init")
	websitesSuper.Dir = "C:/"
	out, err := websitesSuper.CombinedOutput()
	if err != nil {
		updateLogger.Println("Error cloning websites Superproject!!:", err)
		return err
	}

	updateLogger.Println("Cloned websites superproject:", string(out))
	websitesSubmodules.Dir = "C:/websites"
	subOut, err := websitesSubmodules.CombinedOutput()
	if err != nil {
		updateLogger.Panicln("There was a problem initializing subs!!!:", err)
		return err
	}

	updateLogger.Println("Submodule out:", string(subOut))
	return nil

}

func GetWebsiteModule(url string) string {
	updateCmd := exec.Command("git", "submodule", "update", url)
	updateCmd.Dir = "C:/websites"
	out, err := updateCmd.CombinedOutput()
	if err != nil {
		updateLogger.Println("There was an issue updating submodule", url+":", err)
		return string(out)
	}

	return string(out)
}

func pullWebsitesSuperproject() error {
	websitesPull := exec.Command("git", "pull", "http://192.168.1.47:3000/OfflineWebsites/websites.git")
	websitesPull.Dir = "C:/websites"
	out, err := websitesPull.CombinedOutput()
	if err != nil {
		updateLogger.Println("Error pulling websites Superproject!!:", err)
		return err
	}

	updateLogger.Println(string(out))

	return nil

}

func updateClassResources() {
	gitPull := exec.Command("git", "pull", "--force", "origin", "main")
	gitPull.Dir = "./data"
	updateLogger.Println("Checking for class content")
	_, err := gitPull.CombinedOutput()
	if err != nil {
		updateLogger.Println("err:", err)
	}

}
