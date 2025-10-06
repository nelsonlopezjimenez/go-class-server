package CIS

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os/exec"
	"runtime"
)

type FilePath struct {
	Websites   string
	CFC        string
	Videos     string
	ServerPath string
}

//go:embed cis
var build embed.FS

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
		return FilePath{"C:/websites", "C:/Users/Public/CANVAS_FILE_CACHES", "C:/Users/Public/Videos", "C:/Users/Public/classServer"}
	case "darwin":
		return FilePath{"/Users/Shared/websites", "/Users/Shared/CANVAS_FILE_CACHES", "/Users/Shared/Videos", "Users/Shared/ClassServer"}
	default:
		return FilePath{"/var/www/websites", "/var/lib/CANVAS_FILE_CACHES", "/var/lib/Videos", "/var/lib/ClassServer"}
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
		return fmt.Errorf("Unable to launch browser")
	}

	return cmd.Start()
}
