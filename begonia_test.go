package begonia_test

import (
	"testing"
	"github.com/gschier/begonia"
	"github.com/stretchr/testify/assert"
	"os"
	"fmt"
)

func configWithProviders(p ...begonia.Provider) *begonia.Config {
	return &begonia.Config{
		ApplicationConfig: &begonia.ApplicationConfig{Providers: p},
	}
}

type CarServiceInterface interface {
	Honk() string
}

type CarService struct {
	noise string
}

func (cs *CarService) Honk() string {
	if cs.noise == "" {
		return "Honk!"
	}

	return cs.noise
}

type CarServiceProvider struct {
	noise string
}

func (csp *CarServiceProvider) Register(ioc *begonia.Container) {
	ioc.Singleton(func(app *begonia.Application) (*CarService, error) {
		return &CarService{noise: app.Env("honk")}, nil
	})
}

func (csp *CarServiceProvider) Boot(app *begonia.Application) {
	carService := app.Make(new(CarServiceInterface)).(*CarService)
	fmt.Printf("HELLO? %v\n", carService.Honk())
}

func TestApplication_Call(t *testing.T) {
	os.Setenv("honk", "Env Honk!")

	app := begonia.NewApplication(configWithProviders(
		new(CarServiceProvider),
	))

	app.Call(func(c CarServiceInterface) {
		assert.IsType(t, &CarService{}, c, "Should be a car type")
		assert.Equal(t, c.Honk(), "Env Honk!", "Should get honk sound from env")
	})
}

func TestContainer_Make(t *testing.T) {
	app := begonia.NewApplication(configWithProviders())

	app.Bind(func(app *begonia.Application) (*CarService, error) {
		return &CarService{noise: app.Env("honk")}, nil
	})

	instance1 := app.Make(new(CarServiceInterface))
	instance2 := app.Make(new(CarServiceInterface))
	instance3 := app.Make(new(CarService))

	assert.True(t, instance1 != instance2, "Should be new instance")
	assert.True(t, instance1 != instance3, "Should be new instance")
}

func TestContainer_MakeSingleton(t *testing.T) {
	app := begonia.NewApplication(configWithProviders())

	app.Singleton(func(app *begonia.Application) (*CarService, error) {
		return &CarService{noise: app.Env("honk")}, nil
	})

	instance1 := app.Make(new(CarServiceInterface))
	instance2 := app.Make(new(CarServiceInterface))
	instance3 := app.Make(new(CarService))

	assert.True(t, instance1 == instance2, "Should be same instance")
	assert.True(t, instance1 == instance3, "Should be same instance")
}

func TestContainer_MakeInstance(t *testing.T) {
	app := begonia.NewApplication(configWithProviders())

	instance := &CarService{noise: "Honk!"}
	app.Instance(instance)

	instance1 := app.Make(new(CarServiceInterface))
	instance2 := app.Make(new(CarService))

	assert.True(t, instance == instance1, "Should be new instance")
	assert.True(t, instance1 == instance2, "Should be new instance")
}
