package beluga

import (
	"os/exec"
	"strings"
)

func cmdString(cmd string, args ...string) string {
	output, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func gitEnvRead(e Environment) {
	versionTag := cmdString("git", "describe", "--tags", "--match", "v*", "--exclude \"v*-*\"")

	branch := cmdString("git", "rev-parse", "--abbrev-ref", "HEAD")
	environment := ""

	// lists all tags that the current commit points to
	for _, tag := range strings.Split(cmdString("git", "tag", "--points-at", "HEAD"), "\n") {
		if tag == versionTag {
			environment = envProduction
			break
		}
	}
	if environment == "" {
		if branch == e.DefaultBranch() {
			environment = envStaging
		} else {
			environment = envReview
		}
	}

	env := parseVersion(versionTag)
	if len(environment) > 0 && env != nil {
		env[varEnvironment] = environment
	}

	e.MergeMissing(env)
}
