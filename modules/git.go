package CIS

import (
	"bufio"
	"fmt"
	"os"

	command "localhost/CIS/modules/cmd"
)

// GitPull
//
// Creates a command.Command instance of the git pull command and executes it
// updating the repo named in the argument
func GitPull(repoPath string) error {
	pull := command.MakeGitCmd(repoPath, "pull", "--force", "origin", "main")
	return cmdWithOutput(pull)
}

// GitClone
//
// GitClone clones the repo at the url location.
// repoPath: the parent dir containing the repo
// url: repo url
// repoName: name of the repo. This field is required no matter
// the name desired. This allows to check for an existing dir by the
// same name.
func GitClone(repoPath string, url string, repoName string) error {
	_, err := os.Stat(repoPath + "/" + repoName)
	if err == nil {
		return fmt.Errorf("repo already exists")
	}

	clone := command.MakeGitCmd(repoPath, "clone", url, repoName)
	return cmdWithOutput(clone)
}

// cmdWithOutput
//
// This fn takes an instance of command.Command, executes it, and pipes the
// output of stdout and stderr to the server's console.
func cmdWithOutput(c *command.Command) error {
	pipe, err := c.Cmd.StdoutPipe()
	if err != nil {
		return err
	}
	errPipe, err := c.Cmd.StderrPipe()
	if err != nil {
		return err
	}
	err = c.Cmd.Start()
	if err != nil {
		return err
	}

	reader := bufio.NewReader(pipe)
	errRead := bufio.NewReader(errPipe)
	_, err = reader.WriteTo(os.Stdout)
	if err != nil {
		return err
	}

	_, err = errRead.WriteTo(os.Stderr)
	if err != nil {
		return err
	}

	return nil
}
