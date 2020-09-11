package beluga

import (
	"net/url"
	"regexp"
)

// See https://docs.gitlab.com/ee/ci/variables/predefined_variables.html

var gitLabVarMatcher = regexp.MustCompile(`(?i)\$(\$|([a-z0-9_\-]+)|{([a-z0-9_\-]+)})`)

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

	env.MergeMissing(Environment{
		varEnvironment:      environment,
		varRegistry:         e["CI_REGISTRY"],
		varRegistryUsername: e.Get("CI_REGISTRY_USER", "gitlab-ci-token"),
		varRegistryPassword: e["CI_REGISTRY_PASSWORD"],
		varStackName:        e.expand(gitLabVarMatcher, e["CI_PROJECT_PATH_SLUG"]),
		varImagesTemplate:   `{{.Env.CI_REGISTRY_IMAGE}}:{{if .Env.CI_COMMIT_TAG}}{{.Env.CI_COMMIT_TAG}} {{.Env.CI_REGISTRY_IMAGE}}:latest{{else}}{{.Env.CI_COMMIT_REF_NAME}}{{end}}`,
	})

	e.MergeMissing(env)

	// We'll do domain last, so it can include all other vars up to this point.
	domain := gitlabDomain(e)
	e.MergeMissing(Environment{varDomain: domain})

	return nil
}

func gitlabDomain(e Environment) string {
	// If a global variable is used for the job>environment>url field, its
	// variables will not be expanded as expected. This works around that.
	s := e.expand(gitLabVarMatcher, e["CI_ENVIRONMENT_URL"])
	u, err := url.Parse(s)
	if err != nil {
		return ""
	}
	return u.Hostname()
}
