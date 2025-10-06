package CIS

import (
	"embed"
	"io/fs"
	"net/http"
)

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
