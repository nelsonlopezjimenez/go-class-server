// Package logger provides a simple logger that logs messages to the
// the console. It also creates an error log in the same directory as
// the program to log warn and error level messages with stack traces.
package logger

import (
	// "log"

	"fmt"
	"os"
	"runtime"
	"time"
)

// Log Levels:
// [DebugLevel] debugging code and general interest dev msgs
//
// [InfoLevel] General interest msgs for end user
//
// [WarnLevel] Warnings such as possible unexpected behavior or minor errors
//
// [ErrorLevel] serious errors that compromise the functionality of the program
const (
	DebugLevel = 1
	InfoLevel  = 2
	WarnLevel  = 3
	ErrorLevel = 4
)

func Log(level int, message string) {

	if level == 3 || level == 4 {
		preamble := "the following error occurred: "
		message = preamble + message
	}
	logMsg := fmt.Sprintf("[%s] %s: %v", levelToString(level), time.Now().Format(time.DateTime), message)
	fmt.Println(logMsg)
	if level == 1 {
		fmt.Println(GetStackTrace())
	}

	if level == 3 || level == 4 {
		trace := GetStackTrace()

		if level == 4 {
			errorFile, err := os.OpenFile("./errLog.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
			if err != nil {
				fmt.Println(WarnLevel, fmt.Sprintf("The error log could not be opened: %v\n", err))
			}
			defer errorFile.Close()
			fmt.Fprintf(errorFile, "%s\n stack trace: %s\n", logMsg, trace)
		}

		fmt.Fprintf(os.Stderr, "The stack trace for the previous error is: %s", trace)
	}
}

func levelToString(level int) string {
	switch level {
	case 1:
		return "DEBUG"
	case 2:
		return "INFO"
	case 3:
		return "WARN"
	case 4:
		return "ERROR"

	}
	return ""
}

func GetStackTrace() string {
	buf := make([]byte, 1024*16)
	length := runtime.Stack(buf, false)
	return string(buf[:length])
}
