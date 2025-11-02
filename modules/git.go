package CIS

import (
	"bufio"
	"fmt"
	"os"

	command "localhost/CIS/modules/cmd"
)

func GitPull(repoPath string) error {
	pull := command.MakeGitCmd(repoPath, "pull", "--force", "origin", "main")
	return cmdWithOutput(pull)
}
func GitClone(repoPath string, url string, repoName string) error {
	_, err := os.Stat(repoPath+"/"+repoName); if err == nil {
		return fmt.Errorf("repo already exists")
	}

	clone := command.MakeGitCmd(repoPath, "clone", url+repoName )
	return cmdWithOutput(clone)
}

func cmdWithOutput(c *command.Command) error {
pipe, err := c.Cmd.StdoutPipe(); if err != nil {
		return err
	}
	errPipe, err := c.Cmd.StderrPipe(); if err != nil {
		return err
	}
	err = c.Cmd.Start(); if err != nil {
		return err
	}


	reader := bufio.NewReader(pipe)
	errRead := bufio.NewReader(errPipe)
	_, err = reader.WriteTo(os.Stdout); if err != nil {
		return err
	}

	_, err = errRead.WriteTo(os.Stderr); if err != nil {
		return err
	}

	return nil
}