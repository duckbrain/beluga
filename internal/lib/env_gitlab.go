package lib

import (
	"net/url"

	"github.com/gobuffalo/envy"
)

type gitlabEnv struct{}

func (g gitlabEnv) EnvRead() Env {
	if envy.Get("GITLAB_CI", "") == "" {
		return nil
	}

	return Env{
		envDockerUsername: envy.Get("CI_REGISTRY_USER", "gitlab-ci-token"),
		envDockerPassword: envy.Get("CI_REGISTRY_PASSWORD", ""),
		envDockerImage:    envy.Get("CI_REGISTRY_IMAGE", ""),
		envGitRefName:     envy.Get("CI_COMMIT_REF_NAME", ""),
		envDomain:         g.Domain(),
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
