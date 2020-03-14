package beluga

import (
	"net/url"
)

// See https://docs.gitlab.com/ee/ci/variables/predefined_variables.html

func gitlabEnvRead(e Environment) {
	if e["GITLAB_CI"] == "" {
		return
	}

	var environment = envReview
	var image string
	env := parseVersion(e["CI_COMMIT_TAG"])
	if env[varEnvironment] != "" {
		environment = envProduction
	} else if e["CI_COMMIT_REF_NAME"] == e[varGitDefaultBranch] {
		environment = envStaging
	}
	domain := gitlabDomain(e)

	image = e["CI_REGISTRY_IMAGE"] + ":" + e["CI_COMMIT_REF_SLUG"]
	if environment == envStaging {
		image += " " + e["CI_REGISTRY_IMAGE"] + ":latest"
	} else if environment == envProduction {
		image = e["CI_REGISTRY_IMAGE"] + ":" + e[varVersion]
	}

	env.MergeMissing(Environment{
		varEnvironment:      environment,
		varRegistry:         e["CI_REGISTRY"],
		varRegistryUsername: e.Get("CI_REGISTRY_USER", "gitlab-ci-token"),
		varRegistryPassword: e["CI_REGISTRY_PASSWORD"],
		varImage:            image,
		varDomain:           domain,
	})

	e.MergeMissing(env)
}

func gitlabDomain(e Environment) string {
	s := e["CI_ENVIRONMENT_URL"]
	u, err := url.Parse(s)
	if err != nil {
		return ""
	}
	return u.Hostname()
}

func init() {
	envs = append(envs, gitlabEnvRead)
}
