package hemlock_test

import (
	"github.com/gschier/hemlock"
	. "github.com/gschier/hemlock/testutil"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
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
	os.Setenv("honk", "Env Honk!")
	app := NewTestApplication()

	app.Bind(func(app *hemlock.Application) (*CarService, error) {
		return &CarService{Noise: app.Env("honk")}, nil
	})

	instance1 := app.Make(new(CarServiceInterface))
	instance2 := app.Make(new(CarServiceInterface))
	instance3 := app.Make(new(CarService))

	assert.True(t, instance1 != instance2, "Should be new instance")
	assert.True(t, instance1 != instance3, "Should be new instance")
}

func TestContainer_MakeInterface(t *testing.T) {
	os.Setenv("honk", "Env Honk!")
	app := NewTestApplication()

	app.Bind(func(app *hemlock.Application) (CarServiceInterface, error) {
		return &CarService{Noise: app.Env("honk")}, nil
	})

	instance1 := app.Make(new(CarServiceInterface))
	instance2 := app.Make(new(CarServiceInterface))
	instance3 := app.Make(new(CarService))

	assert.True(t, instance1 != instance2, "Should be new instance")
	assert.True(t, instance1 != instance3, "Should be new instance")
}

func TestContainer_MakeInto(t *testing.T) {
	os.Setenv("honk", "Env Honk!")
	app := NewTestApplication()

	app.Bind(func(app *hemlock.Application) (CarServiceInterface, error) {
		return &CarService{Noise: app.Env("honk")}, nil
	})

	var instance1 CarService
	var instance2 CarServiceInterface

	app.MakeInto(&instance1)
	app.MakeInto(&instance2)

	assert.Equal(t, "Env Honk!", instance1.Honk(), "Value should honk")
	assert.Equal(t, "Env Honk!", instance2.Honk(), "Interface should honk")
}

func TestContainer_MakeSingleton(t *testing.T) {
	os.Setenv("honk", "Env Honk!")
	app := NewTestApplication()

	app.Singleton(func(app *hemlock.Application) (*CarService, error) {
		return &CarService{Noise: app.Env("honk")}, nil
	})

	instance1 := app.Make(new(CarServiceInterface))
	instance2 := app.Make(new(CarServiceInterface))
	instance3 := app.Make(new(CarService))

	assert.True(t, instance1 == instance2, "Should be same instance")
	assert.True(t, instance1 == instance3, "Should be same instance")
}

func TestContainer_MakeInstance(t *testing.T) {
	os.Setenv("honk", "Env Honk!")
	app := NewTestApplication()

	instance := &CarService{Noise: "Honk!"}
	app.Instance(instance)

	instance1 := app.Make(new(CarServiceInterface))
	instance2 := app.Make(new(CarService))

	assert.True(t, instance == instance1, "Should be new instance")
	assert.True(t, instance1 == instance2, "Should be new instance")
}
