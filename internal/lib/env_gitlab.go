package lib

import (
	"net/url"
	"strings"

	"github.com/gobuffalo/envy"
)

type gitlabEnv struct{}

func (g gitlabEnv) EnvRead() Env {
	const defaultBranch = "master"

	if envy.Get("GITLAB_CI", "") == "" {
		return nil
	}

	var env = "review"
	var refName = envy.Get("CI_COMMIT_REF_NAME", "")

	if refName == defaultBranch {
		env = "staging"
	}
	if strings.HasPrefix(refName, "v") {
		env = "production"
	}

	return Env{
		varEnvironment:      env,
		varRegistry:         envy.Get("CI_REGISTRY", ""),
		varRegistryUsername: envy.Get("CI_REGISTRY_USER", "gitlab-ci-token"),
		varRegistryPassword: envy.Get("CI_REGISTRY_PASSWORD", ""),
		varImage:            envy.Get("CI_REGISTRY_IMAGE", ""),
		varDomain:           g.Domain(),
	}
}

func (g gitlabEnv) Domain() string {
	e, err := envy.MustGet("CI_ENVIRONMENT_URL")
	if err != nil {
		return ""
	}
	u, err := url.Parse(e)
	if err != nil {
		return ""
	}
	return u.Hostname()
}

func init() {
	envs = append(envs, gitlabEnv{})
}
