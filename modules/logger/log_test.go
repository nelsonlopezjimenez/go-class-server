package logger

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	testCases := []struct {
		level      int
		logMessage string
		want       string
	}{
		{
			InfoLevel,
			"This is an info log",
			fmt.Sprintf("[INFO] %s: This is an info log\n", time.Now().Format(time.DateTime)),
		},
		{
			DebugLevel,
			"This is a debug log",
			fmt.Sprintf("[DEBUG] %s: This is a debug log\n", time.Now().Format(time.DateTime)),
		},
		{
			WarnLevel,
			"This is a warning log",
			fmt.Sprintf("[WARN] %s: the following error occurred: This is a warning log\n", time.Now().Format(time.DateTime)),
		},
		{
			ErrorLevel,
			"This is an error log",
			fmt.Sprintf("[ERROR] %s: the following error occurred: This is an error log\n", time.Now().Format(time.DateTime)),
		},
	}

	for _, testCase := range testCases {
		var originalOutput *os.File
		var isErr bool
		if testCase.level == WarnLevel || testCase.level == ErrorLevel {
			isErr = true
			originalOutput = os.Stderr
		} else {
			originalOutput = os.Stdout
		}
		// Create a pipe to intercept Stdout
		stdRead, stdWrite, err := os.Pipe()
		if err != nil {
			t.Fatalf("Failed to create Stdout pipe:%v", err)
		}
		// Make sure to close the pipe to prevent leaks
		defer func() {
			stdWrite.Close()
			stdRead.Close()
		}()

		// assign Stdout/Stderr to stdWrite
		if isErr {
			os.Stderr = stdWrite
		} else {
			os.Stdout = stdWrite
		}

		//Log test goes here
		Log(testCase.level, testCase.logMessage)

		// Create []byte to hold the contents of stdRead
		readBuffer := make([]byte, 80)
		length, err := stdRead.Read(readBuffer)
		if err != nil {
			t.Fatalf("Failed to read stdRead: %v", err)
		}
		stdRead.Close()
		// read readBuffer
		got := string(readBuffer[:length])
		// reassign Stdout to it's original assignment
		if isErr {
			os.Stderr = originalOutput
		} else {
			os.Stdout = originalOutput
		}

		if got != testCase.want {
			t.Errorf("Got: %q but want: %q", got, testCase.want)
		}
	}

}
