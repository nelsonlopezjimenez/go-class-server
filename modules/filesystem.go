package CIS

import (
	"embed"
	"io/fs"
	"net/http"
)

// The embed directive must point to a valid path relative to this file
// Use all:cis to include all files including those starting with . or _
//
//go:embed all:cis
var cisEmbed embed.FS

// GetIndex
//
// Gets the index.html file embedded into the executable and
// returns the file as a []byte type. This should be called
// by the route handling the "/" route.
func GetIndex() ([]byte, error) {
	index, err := cisEmbed.ReadFile("cis/index.html")
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
	fsys, err := fs.Sub(cisEmbed, "cis/assets")
	if err != nil {
		return nil, err
	}
	return http.FS(fsys), nil
}

// func GetFileSystemHandler() (fs.FS, error) {
// 	func GetFileSystemHandler() (http.FileSystem, error) {
// 			// Try to use embedded filesystem first
// 	fsys, err := fs.Sub(cisEmbed, "cis")
// 	if err != nil {
// 		// If embed fails, fall back to live filesystem
// 		// This allows development without embedded files
// 		if _, statErr := os.Stat("./modules/cis"); statErr == nil {
// 			return os.DirFS("./modules/cis"), nil
// 		}
// 		// Try alternate path for when running from different directory
// 		if _, statErr := os.Stat("./cis"); statErr == nil {
// 			return os.DirFS("./cis"), nil
// 		}
// 		return nil, err
// 	}
// 	return fsys, nil
// }
