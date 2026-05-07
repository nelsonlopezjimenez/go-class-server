package command

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"localhost/CIS/modules/logger"
	"os"
	"os/exec"
)

type CmdArgs string

type Command struct {
	Prog string
	Args []string
	Cmd  *exec.Cmd
}

type CmdError struct {
	err    error
	stderr string
	errMsg string
}

func (e CmdError) Error() string {
	if e.errMsg == "" {
		e.errMsg = fmt.Sprint(e.err, e.stderr)
	}
	return (e.errMsg)
}

func (e CmdError) StdErr() string {
	return e.stderr
}

// MakeGitCmd
//
// Takes supplied args and creates a Command to run a git cmd
func MakeGitCmd(dir string, args ...string) *Command {
	c := &Command{}

	c.Prog = "git"
	c.Args = args
	c.Cmd = exec.Command(c.Prog, c.Args...)
	c.Cmd.Dir = dir
	// _, err := exec.LookPath(c.Prog); if err != nil {
	// 	return nil
	// }
	return c
}

// CheckConfigKeyIsSet
//
// Checks whether the provided git config key is set and
// if it not, sets the key to the value provided
// This is needed to prevent the classServer repo
// failing to commit on pull/merge for new students
// these values can be changed via git cli without issue
func CheckConfigKeyIsSet(key string, value string) error {
	confChk := MakeGitCmd(os.Getenv("PWD"), "config", "get", "--global", key)

	err := confChk.Cmd.Run()
	if err == nil {
		// Key is set and has values
		return nil
	}
	// If the program has gotten to this point, set the
	// key to the supplied values
	// TODO: Specifically check for "exit status 1"
	setKey := MakeGitCmd(os.Getenv("PWD"), "config", "set", "--global", key, value)
	return setKey.Cmd.Run()
}

// Command.ProcessState
//
// returns the state of the corresponding process
func (c Command) ProcessState() string {
	if c.Cmd == nil {
		return ""
	}
	return c.Cmd.ProcessState.String()
}

// Command.AdArgs
//
// Appends the given args to the Command struct
func (c Command) AddArgs(args ...CmdArgs) Command {
	for _, arg := range args {
		c.Args = append(c.Args, string(arg))
	}

	return c
}

// func (c Command) AddOptions(opts ...string) Command {
// 	for _, opt := range opts {
// 		formattedOpt := fmt.Sprint("--opt=")
// 		c.Args = append(c.Args, "--"+opt)
// 	}
// 	return c
// }

// isValidOption
//
// Returns whether arg is valid option or not
func isValidOption(str string) bool {
	return str != "" && str[0] == '-'
}

// CmdStdPipe
//
// Takes a pointer to exec.Cmd and starts the process
// with the process' stdout and stderr piped to the
// main stdout and stderr.
// Note: It is interesting to see what Git outputs on
// the stderr...
func CmdStdPipe(cmd *exec.Cmd) {
	// Create reader for process' stdout
	reader, err := cmd.StdoutPipe()
	if err != nil {
		logger.Log(logger.ErrorLevel, err.Error())
	}

	// Create reader for process' stderr
	errReader, err := cmd.StderrPipe()
	if err != nil {
		logger.Log(logger.ErrorLevel, err.Error())
	}

	// Buffer to store data from the pipes
	buf := make([]byte, 1024)
	errBuf := make([]byte, 1024)

	// Goroutine to handle stdout pipe
	go func() {
		for {

			n, err := reader.Read(buf)
			if err != nil {
				// if there is an error and it is EOF
				// Break out of loop because the no data
				// Will be passed and it is time to end.
				// Otherwise, log the error and notify the user
				if errors.Is(err, io.EOF) || errors.Is(err, fs.ErrClosed) {
					break
				} else {
					logger.Log(logger.ErrorLevel, fmt.Sprintf("From Stdout: %v, %T", err.Error(), err))
				}
			}
			if n == 0 {
				break
			}
			logger.Log(logger.InfoLevel, string(buf[:n]))
		}
	}()

	// Goroutine to handle stdout pipe
	go func() {
		for {
			n, err := errReader.Read(errBuf)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				} else {
					logger.Log(logger.ErrorLevel, fmt.Sprintf("From Stderr: %v", err.Error()))
				}
			}
			if n == 0 {
				break
			}
			logger.Log(logger.WarnLevel, string(errBuf[:n]))

		}
	}()

	//  Start the process
	cmd.Start()
	// Wait for the process to end
	cmd.Wait()
	// Close the read streams
	reader.Close()
	errReader.Close()
}
