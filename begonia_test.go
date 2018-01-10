package begonia_test

import (
	"testing"
	"github.com/gschier/begonia"
	"github.com/stretchr/testify/assert"
	"os"
	"fmt"
)

func TestApplication_Call(t *testing.T) {
	os.Setenv("honk", "Env Honk!")

	app := NewTestApplication(new(CarServiceProvider))

	app.Call(func(c CarServiceInterface) {
		assert.IsType(t, &CarService{}, c, "Should be a car type")
		assert.Equal(t, c.Honk(), "Env Honk!", "Should get honk sound from env")
	})
}

func TestContainer_Make(t *testing.T) {
	app := NewTestApplication()

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
	app := NewTestApplication()

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
	app := NewTestApplication()

	instance := &CarService{noise: "Honk!"}
	app.Instance(instance)

	instance1 := app.Make(new(CarServiceInterface))
	instance2 := app.Make(new(CarService))

	assert.True(t, instance == instance1, "Should be new instance")
	assert.True(t, instance1 == instance2, "Should be new instance")
}
