// DO NOT EDIT: This file is generated by var_gen.go

package beluga

const (
	varApplication      = "BELUGA_APPLICATION"
	varDockerContext    = "BELUGA_DOCKER_CONTEXT"
	varDockerfile       = "BELUGA_DOCKERFILE"
	varDomain           = "BELUGA_DOMAIN"
	varEnvironment      = "BELUGA_ENVIRONMENT"
	varGitDefaultBranch = "BELUGA_GIT_DEFAULT_BRANCH"
	varImage            = "BELUGA_IMAGE"
	varRegistry         = "BELUGA_REGISTRY"
	varRegistryPassword = "BELUGA_REGISTRY_PASSWORD"
	varRegistryUsername = "BELUGA_REGISTRY_USERNAME"
	varStackName        = "BELUGA_STACK_NAME"
	varVariant          = "BELUGA_VARIANT"
	varVersion          = "BELUGA_VERSION"
)

var knownVarNames = []string{
	varApplication,
	varDockerContext,
	varDockerfile,
	varDomain,
	varEnvironment,
	varGitDefaultBranch,
	varImage,
	varRegistry,
	varRegistryPassword,
	varRegistryUsername,
	varStackName,
	varVariant,
	varVersion,
}

// Context of the docker build, defaults to root of the project
func (e Environment) DockerContext() string {
	return e[varDockerContext]
}

// Dockerfile to use in docker build, defaults `Dockerfile` in the context
// directory (like docker does)
func (e Environment) Dockerfile() string {
	return e[varDockerfile]
}

// Domain name to deploy the stack to. This will be passed to the environment
// when doing the docker deploy, so the compose file can reference this
// appropriately.
func (e Environment) Domain() string {
	return e[varDomain]
}

// Docker image path to push to after build
func (e Environment) Image() string {
	return e[varImage]
}

// Docker registry to log into before pushing
func (e Environment) Registry() string {
	return e[varRegistry]
}

// Password to use to log into Docker registry
func (e Environment) RegistryPassword() string {
	return e[varRegistryPassword]
}

// Username to use to log into Docker registry
func (e Environment) RegistryUsername() string {
	return e[varRegistryUsername]
}

func (e Environment) StackName() string {
	return e[varStackName]
}
