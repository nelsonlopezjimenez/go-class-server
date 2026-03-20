package util

import (
	"os"
	"testing"
)

func TestCreateSymlink(t *testing.T) {
	linksToCreate := []string{
		"CANVAS_FILE_CACHES",
		"websites",
		"Videos",
		"serverPath",
		"static",
		"Public",
	}
	FilePaths := GetOSPaths()
	for _, v := range linksToCreate {
		err := CreateSymlink(FilePaths[v], os.TempDir()+"/"+v)
		if err != nil {
			t.Fatalf("Expected nil but got %v", err)
		}

		defer os.Remove(os.TempDir() + "/" + v)
	}
}
