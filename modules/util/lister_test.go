package util

import (
	"fmt"
	"os"
	"testing"
)

func TestListBuilder(t *testing.T) {
	testDir, err := os.MkdirTemp("", "test_dir")
	defer os.RemoveAll(testDir)
	if err != nil {
		t.Fatalf("Failed to create tmp dir for ListBuilder: %v", err)
	}
	for i := range 10 {
		err := os.Mkdir(fmt.Sprintf("%s/testFolder%d", testDir, i), 0644)
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
	}
	testRoot := MakeRootDir(testDir)
	testCases := struct {
		dirList map[string]bool
		want    map[string]bool
	}{testRoot.ListBuilder(), map[string]bool{"testFolder0": true, "testFolder1": true, "testFolder2": true, "testFolder3": true, "testFolder4": true, "testFolder5": true, "testFolder6": true, "testFolder7": true, "testFolder8": true, "testFolder9": true}}
	// sort.Strings(testCases.dirList)
	// sort.Strings(testCases.want)
	if len(testCases.dirList) != len(testCases.want) {
		t.Errorf("Maps are not the same length. Expected %v, but got %v", testCases.want, testCases.dirList)
	}
	for key, _ := range testCases.want {
		if !testCases.dirList[key] {
			t.Errorf("Expected %v, but got %v", testCases.want, testCases.dirList)
		}
	}
}

func TestFindIndex(t *testing.T) {
	testDir, err := os.MkdirTemp("", "test_dir")
	if err != nil {
		t.Fatalf("Failed to create tmp dir for FindIndex: %v", err)
	}
	defer os.RemoveAll(testDir)
	for i := range 6 {
		err := os.Mkdir(fmt.Sprintf("%s/testFolder%d", testDir, i), 0644)
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		if i%2 != 0 {
			err := os.WriteFile(fmt.Sprintf("%s/testFolder%d/not_index.html", testDir, i), []byte("This is an html file"), 0644)
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			continue
		}
		err = os.WriteFile(fmt.Sprintf("%s/testFolder%d/index.html", testDir, i), []byte("This is an html file"), 0644)
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
	}
	testCases := []struct {
		got  bool
		want bool
	}{
		{
			MakeRootDir(testDir + "/testFolder0").HasIndex(),
			true,
		},
		{
			MakeRootDir(testDir + "/testFolder1").HasIndex(),
			false,
		},
		{
			MakeRootDir(testDir + "/testFolder2").HasIndex(),
			true,
		},
		{
			MakeRootDir(testDir + "/testFolder3").HasIndex(),
			false,
		},
		{
			MakeRootDir(testDir + "/testFolder4").HasIndex(),
			true,
		},
		{
			MakeRootDir(testDir + "/testFolder5").HasIndex(),
			false,
		},
	}
	for _, tc := range testCases {
		if tc.got != tc.want {
			t.Errorf("Got %v but wanted %v", tc.got, tc.want)
		}
	}
}
