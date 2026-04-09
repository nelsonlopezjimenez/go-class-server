// Package logger provides a simple logger that logs messages to the
// the console. It also creates an error log in the same directory as
// the program to log warn and error level messages with stack traces.
package logger

// TODO: Create a logger struct with methods instead of stand alone fns

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
	logColor, logLevel := levelToString(level)

	logMsg := fmt.Sprintf("%s[%s] %s: %v\033[0m", logColor, logLevel, time.Now().Format(time.DateTime), message)

	if level == 1 {
		message += GetStackTrace()
	}

	if level == 1 || level == 2 {
		fmt.Println(logMsg)
	}

	if level == 3 || level == 4 {
		trace := GetStackTrace()

		if level == 4 {
			printErrToFile(logMsg, trace)
			fmt.Fprintf(os.Stderr, "The stack trace for the previous error is: %s", trace)
		}

		fmt.Fprintf(os.Stderr, "%s\n", logMsg)
	}
}

func levelToString(level int) (string, string) {
	switch level {
	case 1:
		return "\033[36m", "DEBUG"
	case 2:
		return "\033[32m", "INFO"
	case 3:
		return "\033[33m", "WARN"
	case 4:
		return "\033[31m", "ERROR"

	}
	return "", ""
}

// GetStackTrace
//
// Prints the stack trace of the current process. Meant to give Log()
// more verbose error and/or debug messages.
func GetStackTrace() string {
	buf := make([]byte, 1024*16)
	length := runtime.Stack(buf, false)
	return string(buf[:length])
}

func printErrToFile(logMsg string, trace string) {
	errorFile, err := os.OpenFile("./errLog.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println(WarnLevel, fmt.Sprintf("The error log could not be opened: %v\n", err))
	}
	defer errorFile.Close()
	fmt.Fprintf(errorFile, "%s\n stack trace: %s\n", logMsg, trace)

}
