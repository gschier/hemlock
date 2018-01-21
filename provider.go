package hemlock

import (
	"github.com/gschier/hemlock/interfaces"
)

type Provider interface {
	// Register registers a new provider. Any setup should happen here
	Register(interfaces.Container)

	// Boot is called after all service providers have been registered
	Boot(*Application) error
}
