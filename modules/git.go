package CIS

import (
	"fmt"
	"os"

	command "localhost/CIS/modules/cmd"
)

// GitPull
//
// Creates a command.Command instance of the git pull command and executes it
// updating the repo named in the argument
func GitPull(repoPath string) ([]byte, error) {
	pull := command.MakeGitCmd(repoPath, "pull", "--force", "origin", "main")
	// return cmdWithOutput(pull)
	return pull.Cmd.CombinedOutput()
}

// GitClone
//
// GitClone clones the repo at the url location.
// repoPath: the parent dir containing the repo
// url: repo url
// repoName: name of the repo. This field is required no matter
// the name desired. This allows to check for an existing dir by the
// same name.
func GitClone(repoPath string, url string, repoName string) ([]byte, error) {
	_, err := os.Stat(repoPath + "/" + repoName)
	if err == nil {
		return nil, fmt.Errorf("repo already exists")
	}

	clone := command.MakeGitCmd(repoPath, "clone", url)

	// return cmdWithOutput(clone)
	return clone.Cmd.CombinedOutput()
}
