package external

import (
	"fmt"
	"os"

	"localhost/CIS/modules/command"
)

// GitPull
//
// Creates a command.Command instance of the git pull command and executes it
// updating the repo named in the argument
func GitPull(repoPath string) {
	pull := command.MakeGitCmd(repoPath, "pull", "--force", "origin", "main")
	// return cmdWithOutput(pull)
	command.CmdStdPipe(pull.Cmd)
	// return pull.Cmd.CombinedOutput()
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

	clone := command.MakeGitCmd(repoPath, "clone", url)
	command.CmdStdPipe(clone.Cmd)

	// return cmdWithOutput(clone)
	return nil
}
