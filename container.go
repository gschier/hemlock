package hemlock

type Container interface {
	// Bind binds the type of v as a dependency
	Bind(fn interface{})

	// Singleton binds the type of v as a dependency. Will only get instantiated once
	Singleton(fn interface{})

	// Instance binds an already-created value as a dependency
	Instance(i interface{})
}
