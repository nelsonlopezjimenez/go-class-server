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
	// return cmdWithOutput(pull)
	return pull.Cmd.Run()
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
	fmt.Println("clone stat")
	_, err := os.Stat(repoPath + "/" + repoName)
	if err == nil {
		return fmt.Errorf("repo already exists")
	}
	fmt.Println("clone mkgit")

	clone := command.MakeGitCmd(repoPath, "clone", url)
	fmt.Println("cmdwith")

	// return cmdWithOutput(clone)
	return clone.Cmd.Run()
}

// cmdWithOutput
//
// This fn takes an instance of command.Command, executes it, and pipes the
// output of stdout and stderr to the server's console.
func cmdWithOutput(c *command.Command) error {
	// TODO: fn needs to not block at the reader. This needs to be addressed before this fn can be reimplemented.

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

	outBuf := make([]byte, 1024)
	reader := bufio.NewReader(pipe)
	n, err := reader.Read(outBuf)
	if err != nil {
		return fmt.Errorf("err with outbuf: %w", err)
	}
	_, err = bufio.NewWriter(os.Stdout).Write(outBuf[:n])
	if err != nil {
		return fmt.Errorf("err with outbuf: %w", err)
	}
	err = bufio.NewWriter(os.Stdout).Flush()
	if err != nil {
		return err
	}

	errRead := bufio.NewReader(errPipe)
	if _, err := bufio.NewWriter(os.Stdout).Write(outBuf[:n]); err != nil {
		return err
	}

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
