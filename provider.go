package hemlock

type Providers []Provider

type Provider interface {
	// Register registers a new provider. Any setup should happen here
	Register(Container)

	// Boot is called after all service providers have been registered
	Boot(*Application)
}
