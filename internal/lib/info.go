package lib

import (
	"net/url"

	"github.com/gobuffalo/envy"
)

type Info interface {
	IsApplicable() bool
	DockerUsername() string
	DockerPassword() string
	DockerImage() string
	Domain() string
}

type BelugaInfo struct{}

func (i BelugaInfo) IsApplicable() bool {
	return true
}
func (i BelugaInfo) DockerUsername() string {
	return envy.Get("BELUGA_DOCKER_USERNAME", "")
}
func (i BelugaInfo) DockerPassword() string {
	return envy.Get("BELUGA_DOCKER_PASSWORD", "")
}
func (i BelugaInfo) DockerImage() string {
	return envy.Get("BELUGA_DOCKER_IMAGE", "")
}
func (i BelugaInfo) Domain() string {
	return envy.Get("BELUGA_DOMAIN", "")
}

type GitlabInfo struct{}

func (i GitlabInfo) IsApplicable() bool {
	return len(envy.Get("GITLAB_CI", "")) > 0
}
func (i GitlabInfo) DockerUsername() string {
	return envy.Get("CI_REGISTRY_USER", "gitlab-ci-token")
}
func (i GitlabInfo) DockerPassword() string {
	return envy.Get("CI_REGISTRY_PASSWORD", "")
}
func (i GitlabInfo) DockerImage() string {
	return envy.Get("CI_REGISTRY_IMAGE", "")
}
func (i GitlabInfo) Domain() string {
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
