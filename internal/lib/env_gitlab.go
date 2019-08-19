package lib

import (
	"net/url"
	"strings"
)

// See https://docs.gitlab.com/ee/ci/variables/predefined_variables.html

type gitlabEnv struct{}

func (g gitlabEnv) EnvRead(e Environment) {
	const defaultBranch = "master"

	if e["GITLAB_CI"] == "" {
		return
	}

	var env = envReview
	if strings.HasPrefix(e["CI_COMMIT_TAG_NAME"], "v") {
		env = envProduction
	} else if e["CI_COMMIT_REF_NAME"] == e.GitDefaultBranch() {
		env = envStaging
	}

	e.MergeMissing(Environment{
		varEnvironment:      env,
		varRegistry:         e["CI_REGISTRY"],
		varRegistryUsername: e.Get("CI_REGISTRY_USER", "gitlab-ci-token"),
		varRegistryPassword: e["CI_REGISTRY_PASSWORD"],
		varImage:            e["CI_REGISTRY_IMAGE"],
		varDomain:           g.Domain(e),
	})
}

func (g gitlabEnv) Domain(e Environment) string {
	s := e["CI_ENVIRONMENT_URL"]
	u, err := url.Parse(s)
	if err != nil {
		return ""
	}
	return u.Hostname()
}

func init() {
	envs = append(envs, gitlabEnv{})
}
