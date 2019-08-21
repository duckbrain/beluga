package lib

import (
	"os/exec"
	"strings"
)

func cmdString(cmd string, args ...string) string {
	output, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return ""
	}
	return string(output)
}

func gitEnvRead(e Environment) {
	version := parseVersion(cmdString("git", "describe", "--tags", "--match", "v*"))

	branch := cmdString("git", "rev-parse", "--abbrev-ref", "HEAD")
	environment := ""

	// lists all tags that the current commit points to
	for _, tag := range strings.Split(cmdString("git", "tag", "--points-at", "HEAD"), "\n") {
		if tag == version {
			environment = envProduction
			break
		}
	}
	if environment == "" {
		if branch == e.GitDefaultBranch() {
			environment = envStaging
		} else {
			environment = envReview
		}
	}

	e.MergeMissing(Environment{
		varVersion:     version,
		varEnvironment: environment,
	})
}
