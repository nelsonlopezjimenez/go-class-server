package util

import (
	"embed"
	"fmt"
	"io/fs"
	"localhost/CIS/modules/logger"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

type FilePath = map[string]string

//go:embed cis
var build embed.FS
var FilePathMap FilePath

// GetIndex
//
// Gets the index.html file embedded into the executable and
// returns the file as a []byte type. This should be called
// by the route handling the "/" route.
func GetIndex() ([]byte, error) {
	index, err := build.ReadFile("cis/index.html")
	if err != nil {
		return nil, err
	}
	return index, nil
}

// GetFileSystemHandler
//
// Returns type http.FileSystem to serve embedded assets for
// the frontend. This should be called and implemented by
// StaticFS.
func GetFileSystemHandler() (http.FileSystem, error) {
	fsys, err := fs.Sub(build, "cis/assets")
	if err != nil {
		return nil, err
	}
	return http.FS(fsys), nil
}

// GetOSPaths
// Returns a struct with the appropriate paths depending on OS
func GetOSPaths() FilePath {
	switch runtime.GOOS {
	case "windows":
		return FilePath{
			"websites":           "C:/websites",
			"CANVAS_FILE_CACHES": "C:/Users/Public/CANVAS_FILE_CACHES",
			"Videos":             "C:/Users/Public/Videos",
			"serverPath":         "C:/Users/Public/classServer",
			"Public":             "C:/Users/Public",
			"static":             "C:/Users/Public/classServer/static",
		}
	case "darwin":
		return FilePath{
			"websites":           "/Users/Shared/websites",
			"CANVAS_FILE_CACHES": "/Users/Shared/CANVAS_FILE_CACHES",
			"Videos":             "/Users/Shared/Videos",
			"serverPath":         "Users/Shared/ClassServer",
			"Public":             "Users/Shared",
			"static":             "Users/Shared/ClassServer/static",
		}
	default:
		usr, _ := os.UserHomeDir()
		return FilePath{
			"websites":           usr + "/websites",
			"CANVAS_FILE_CACHES": usr + "/CANVAS_FILE_CACHES",
			"Videos":             usr + "/Videos",
			"serverPath":         usr + "/classServer",
			"Public":             usr,
			"static":             usr + "classServer/static",
		}
	}
}

// OpenBrowser opens the server frontend in the OS dependent browser
func OpenBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		browsers := []string{"xdg-open", "google-chrome", "firefox", "chromium"}
		for _, browser := range browsers {
			cmd = exec.Command(browser, url)
			break
		}
	}
	if cmd == nil {
		return fmt.Errorf("unable to launch browser")
	}

	return cmd.Start()
}

func GetRoot() string {
	switch runtime.GOOS {
	case "windows":
		return "C:/"
	case "darwin":
		return "Users/Shared"
	default:
		usr, _ := os.UserHomeDir()

		return usr
	}
}

func DeleteSite(dir string) error {
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}
	return nil
}

func GetDepPath() string {
	var depPath string

	cwd, err := os.Getwd()
	if err != nil {
		// fmt.Println("Cannot get CWD:", err)
		logger.Log(logger.WarnLevel, fmt.Sprintf("Unable to get CWD: %v. Some features may not work properly.", err))
	}
	_, dirErr := os.Stat(cwd + "/data")

	if dirErr == nil {
		depPath = cwd + "/data"
	} else {
		depPath = cwd + "/ClassServerResources"
	}

	return depPath
}

func CreateSymlink(linkName string, linkTarget string) error {

	err := os.Symlink(linkName, linkTarget)
	if err != nil {
		return fmt.Errorf("there was an error creating the link to %s: %w", linkName, err)
	}
	return nil
}

func CheckForSymlink(linkName string) error {
	osPaths := GetOSPaths()
	_, err := os.Stat(osPaths["static"] + "/" + linkName)
	if err != nil {
		if os.IsNotExist(err) {
			return CreateSymlink(osPaths[linkName], osPaths["static"]+"/"+linkName)
		} else {
			return fmt.Errorf("could not check for link %s: %w", linkName, err)
		}
	}
	return nil
}

func IsWindows() bool {
	return runtime.GOOS == "windows"
}
