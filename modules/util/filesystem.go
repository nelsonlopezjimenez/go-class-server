package util

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

type FilePath struct {
	Websites   string
	CFC        string
	Videos     string
	ServerPath string
}

//go:embed cis scripts
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
		usr, _ := os.UserHomeDir()
		return FilePath{usr + "/websites", usr + "/CANVAS_FILE_CACHES", usr + "/Videos", usr + "/classServer"}
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

func RunMigrateScript() error {
	tmpDir, err := os.MkdirTemp(os.TempDir(), "cisMigration-")
	if err != nil {
		return err
	}

	tmpFile, err := os.OpenFile(tmpDir+"/migrate.sh", os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}

	scripts, _ := build.ReadFile("scripts/migrate.sh")

	_, err = tmpFile.Write(scripts)
	if err != nil {
		return err
	}
	tmpFile.Close()

	migrate := exec.Command("C:/Program Files/Git/git-bash.exe", tmpFile.Name())
	migrate.Dir = GetOSPaths().Websites
	err = migrate.Run()
	if err != nil {
		fmt.Println(err)
	}

	err = os.RemoveAll(tmpDir)
	if err != nil {
		return err
	}

	return nil
}
