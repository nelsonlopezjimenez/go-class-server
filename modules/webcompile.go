// Rocky Connor 420711
package CIS

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunCode takes the  supplied code and code language and executes it.
func RunCode(code string, lang string) (_ []byte, err error) {
	// defer will wait to run the specified code until the fn finishes execution
	defer func() {
		if err != nil {
			err = fmt.Errorf("in code compiler: %w", err)
		}
	}()
	// Declare and initialize some vars for use later on
	ext := ""
	cmdSlice := []string{}
	// Conditional to set the vars we declared previously depending on language value
	tmpDir := strings.ReplaceAll(os.TempDir(), "\\", "/")
	switch lang {
	case "go":
		ext = "go"
		cmdSlice = append(cmdSlice, "go", "run")
	case "python":
		ext = "py"
		cmdSlice = append(cmdSlice, "python")
	case "perl":
		ext = "pl"
		cmdSlice = append(cmdSlice, "perl")
	case "typescript":
		_, checkErr := os.ReadDir(tmpDir + "/node_modules")
		if os.IsNotExist(checkErr) {
			npm := exec.Command("npm", "install", "typescript")
			npm.Dir = tmpDir
			npm.Run()
		}
		ext = "ts"
		cmdSlice = append(cmdSlice, "npx", "tsc")
	}

	// Here we must unescape the sanitized form data from the w3schools examples page
	unEscapedCode := strings.ReplaceAll(code, "w3equalsign ", "=")
	unEscapedCode = strings.ReplaceAll(unEscapedCode, "w3plussign", "+")

	// Create a tmp file to store the code to be executed
	tmpFile, err := os.CreateTemp("", ext+"*."+ext)
	if err != nil {
		return nil, err
	}
	// fileS contains a string representation of the temp file path
	fileS := strings.ReplaceAll(tmpFile.Name(), "\\", "/")
	// Write the unescaped code to the tmp file
	cmdSlice = append(cmdSlice, fileS)
	tmpFile.WriteString(unEscapedCode)
	// execute the code using the arguments supplied by the switch fn
	echoCmd := mkCmd(cmdSlice...)
	echoCmd.Dir = tmpDir
	// Code result
	out, err := echoCmd.CombinedOutput()
	if err != nil {
		fmt.Println("echoCmd errored", tmpFile.Name())
		return out, err
	}

	// Once finished with the tmp file, close it
	tmpFile.Close()
	jsFile := strings.Replace(tmpFile.Name(), ".ts", ".js", 1)
	if lang == "typescript" {
		tsRun := mkCmd("node", jsFile)
		tsRun.Dir = tmpDir
		jsOut, err := tsRun.CombinedOutput()
		if err != nil {
			return jsOut, err
		}
		out = jsOut
		os.Remove(jsFile)
	}

	// Delete the tmp file now that we have finished with it
	tmpErr := os.Remove(fileS)
	if tmpErr != nil {
		fmt.Println("There was an issue deleting the tmp code file:", tmpErr)
	}

	return out, nil
}

func mkCmd(cmds ...string) *exec.Cmd {
	prog := cmds[0]
	argsEnd := len(cmds)
	argSlice := cmds[1:argsEnd]
	return exec.Command(prog, argSlice...)
}
