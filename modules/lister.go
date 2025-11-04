// Rocky Connor 420711
package CIS

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type RootDir struct {
	Root string
}

type LinkList struct {
	Title []string
	Path  []string
}

type Meta struct {
	Week        int
	Title       string
	Description string
}

type LessonInfo struct {
	Week        int       `json:"week"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	FileSize    int64     `json:"file_size"`
	Section     string    `json:"section"`
	Content     string    `json:"content"`
}

// FindIndex
// Changed name from Recursive to FindIndex in order to better explain
// the fn's purpose.
//
// This function will recursively search the given dir looking specifically
// for the first index.html file in its subdirectories and pass any found
// to a callback for user defined processing.
func (dir RootDir) FindIndex(cb func(path string, ent fs.DirEntry)) {
	root := os.DirFS(dir.Root)
	rootDir, err := fs.ReadDir(root, ".")
	if err != nil {
		fmt.Println("There was an error opening the dir", err)
	}

	var haveIndex bool = false
	for _, entry := range rootDir {
		if entry.Type().IsRegular() && entry.Name() == "index.html" {
			haveIndex = true

			cb(dir.Root, entry)
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

	for _, ent := range list {
		if ent.IsDir() && !strings.HasPrefix(ent.Name(), ".") {
			linkList = append(linkList, ent.Name())
		}
	}

	return linkList
}

//	(RootDir).HasIndex
//
// This method recursively searches the supplied path
// and returns true if the directory or a child contains
// an index.html file. If it does not, the method returns false
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

// MakeLessonInfo
//
// Creates and returns a struct with the data for the specified markdown file.
// It takes MD files with front matter and parses it into a go struct that can
// then be used as needed. The struct is typically sent to the client at a json obj.
func MakeLessonInfo(dir string, lessonFile string, serverLog *log.Logger) LessonInfo {
	var lesson LessonInfo
	var metaData Meta

	// creates a fs.FS  for the information directory
	fsys := os.DirFS("./data/markdown/lessons/" + dir)
	// Opens the requested markdown file
	// TODO: Should be after os.Stat
	file, err := fs.ReadFile(fsys, lessonFile)
	if err != nil {
		serverLog.Panicln("There was an error getting the requested file:", err)
	}

	// Checks to make sure the file exists. If it does, it populates the CreatedAt, FileSize,
	// and Section fields.
	fileInfo, err := os.Stat("./data/markdown/lessons/" + dir + "/" + lessonFile)
	if err != nil {
		// TODO: Needs to return err and have err handling rather than just swallowing.
		serverLog.Println(err)
	} else {
		lesson.CreatedAt = fileInfo.ModTime()
		lesson.FileSize = fileInfo.Size()
		lesson.Section = dir
	}

	// Converts the returned []byte into a string for manipulation
	contentStr := string(file)

	// Looks for the front matter fences and parses out the keys inside
	// the fence as yaml and adds them to the LessonInfo struct
	if strings.HasPrefix(contentStr, "---") {
		contentSlice := strings.SplitN(contentStr, "---", 3)
		err = yaml.Unmarshal([]byte(contentSlice[1]), &metaData)
		if err != nil {
			// TODO: Handle the error better
			fmt.Println(err)
		} else {
			lesson.Title = metaData.Title
			lesson.Week = metaData.Week
			lesson.Description = metaData.Description
			lesson.Content = strings.TrimSpace(contentSlice[2])
		}
	}
	// Return the struct ready for use elsewhere.
	return lesson
}

func MakeRootDir(dir string) RootDir {
	return RootDir{Root: dir}
}
