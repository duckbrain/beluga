package beluga

import (
	"net/url"
)

// See https://docs.gitlab.com/ee/ci/variables/predefined_variables.html

func gitlabEnvRead(e Environment) error {
	if e["GITLAB_CI"] == "" {
		return nil
	}

	if e[varDefaultBranch] == "" {
		e[varDefaultBranch] = e["CI_DEFAULT_BRANCH"]
	}

	var environment = envReview
	env := parseVersion(e["CI_COMMIT_TAG"])
	if env[varEnvironment] != "" {
		environment = envProduction
	} else if e["CI_COMMIT_REF_NAME"] == e[varDefaultBranch] {
		environment = envStaging
	}
	domain := gitlabDomain(e)

	env.MergeMissing(Environment{
		varEnvironment:      environment,
		varRegistry:         e["CI_REGISTRY"],
		varRegistryUsername: e.Get("CI_REGISTRY_USER", "gitlab-ci-token"),
		varRegistryPassword: e["CI_REGISTRY_PASSWORD"],
		varDomain:           domain,
		varStackName:        e["CI_PROJECT_PATH_SLUG"],
		varImagesTemplate:   `{{.Env.CI_REGISTRY_IMAGE}}:{{if .Env.CI_COMMIT_TAG}}{{.Env.CI_COMMIT_TAG}} {{.Env.CI_REGISTRY_IMAGE}}:latest{{else}}{{.Env.CI_COMMIT_REF_NAME}}{{end}}`,
	})

	e.MergeMissing(env)

	return nil
}

func gitlabDomain(e Environment) string {
	s := e["CI_ENVIRONMENT_URL"]
	u, err := url.Parse(s)
	if err != nil {
		return ""
	}
	return u.Hostname()
}
