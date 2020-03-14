package beluga

type DeployMode int

const (
	ComposeMode DeployMode = 2
	SwarmMode   DeployMode = 1
)

type Deployer interface {
	Deploy(opts DeployOpts) error
	Teardown(opts DeployOpts) error
}

// func (env Environment) Logger() logrus.StdLogger {
// 	logger := logrus.New()
// 	logger.SetLevel(logrus.DebugLevel)
// 	return logger
// }

// func (env Environment) Deployer() Deployer {
// 	logger := env.Logger()
// 	host := env[varDeployDockerHost]

// 	switch host[:strings.Index(host, ":")] {
// 	case "portainer", "portainer-insecure":
// 		logger.Printf("Portainer url: %v", host)
// 		deployer := &PortainerDeploy{}
// 		client, err := portainer.New(host, deployer)
// 		if err != nil {
// 			panic(err)
// 		}
// 		client.Logger = logger
// 		deployer.Client = client
// 		return deployer
// 	case "ssh":
// 		panic("SSH not implemented")
// 	default:
// 		return Docker(host)
// 	}
// }
