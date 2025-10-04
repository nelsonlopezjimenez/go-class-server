// Rocky Connor 420711
package CIS

import (
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"strings"
)

type RootDir struct {
	Root string
}

type LinkList struct {
	Title []string
	Path  []string
}

// Changed name from Recursive to FindIndex in order to better explain
// the fn's purpose.
//
// This function will recursively search the given dir looking specifically
// for the first index.html file in its subdirectories and pass any found
// to a callback for user defined processing.
func (dir RootDir) FindIndex(cb func(path string, fileName string)) {
	root := os.DirFS(dir.Root)
	rootDir, err := fs.ReadDir(root, ".")
	if err != nil {
		fmt.Println("There was an error opening the dir", err)
	}

	var haveIndex bool = false
	for _, entry := range rootDir {
		if entry.Type().IsRegular() && entry.Name() == "index.html" {
			haveIndex = true

			cb(dir.Root, entry.Name())
		}
	}

	if !haveIndex {
		for _, entry := range rootDir {
			if !entry.Type().IsRegular() && !haveIndex {
				deeperLook := RootDir{dir.Root + "/" + entry.Name()}
				deeperLook.FindIndex(cb)

			}
		}
	}
}

func (dir RootDir) RecursiveSearch(ext string, cb func(path string, fileName string)) {
	root := os.DirFS(dir.Root)
	rootDir, err := fs.ReadDir(root, ".")
	if err != nil {
		fmt.Println("There was an error opening the dir", err)
	}

	for _, entry := range rootDir {
		deeperLook := RootDir{dir.Root + "/" + entry.Name()}
		if !entry.Type().IsRegular() {
			deeperLook.RecursiveSearch(ext, cb)
		}
		if findFileExt(entry.Name(), ext) && entry.Type().IsRegular() {
			currentDir := strings.Split(dir.Root, "/")
			cb(currentDir[len(currentDir)-1], entry.Name())
		}
	}
}

// This fn returns an array of strings representing the contents of the RootDir passed to it
func (dir RootDir) ListBuilder() []string {
	fsys := os.DirFS(dir.Root)
	linkList := []string{}
	list, err := fs.ReadDir(fsys, ".")
	if err != nil {
		fmt.Println(err)
	}

	for i := range list {
		if list[i].IsDir() {
			linkList = append(linkList, list[i].Name())
		}
	}

	return linkList
}

func (dir RootDir) HasIndex() bool {
	root := os.DirFS(dir.Root)
	rootDir, err := fs.ReadDir(root, ".")
	if err != nil {
		fmt.Println("There was an error opening the dir", err)
	}

	var haveIndex bool = false
	for _, entry := range rootDir {
		if entry.Type().IsRegular() && entry.Name() == "index.html" {
			haveIndex = true

			return haveIndex
		}
	}

	if !haveIndex {
		for _, entry := range rootDir {
			if !entry.Type().IsRegular() && !haveIndex {
				deeperLook := RootDir{dir.Root + "/" + entry.Name()}
				deeperLook.HasIndex()

			}
		}
	}
	return haveIndex
}

func findFileExt(name string, ext string) bool {
	re := regexp.MustCompile(ext + "$")

	return re.MatchString(name)

}
