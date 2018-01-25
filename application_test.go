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
	app := NewTestApplication(
		new(CarServiceProvider),
		new(StringServiceProvider),
	)

	app.ResolveInto(func(c CarServiceInterface, s string) {
		assert.IsType(t, &CarService{}, c, "Should be a car type")
		assert.Equal(t, c.Honk(), "Env Honk!", "Should get honk sound from env")
		assert.Equal(t, s, "Hello!", "Should get string")
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

func TestContainer_Resolve(t *testing.T) {
	os.Setenv("honk", "Env Honk!")
	app := NewTestApplication()

	app.Bind(func(app *hemlock.Application) (CarServiceInterface, error) {
		return &CarService{Noise: app.Env("honk")}, nil
	})

	var instance1 CarService
	var instance2 CarServiceInterface

	app.Resolve(&instance1)
	app.Resolve(&instance2)

	assert.Equal(t, "Env Honk!", instance1.Honk(), "Value should honk")
	assert.Equal(t, "Env Honk!", instance2.Honk(), "Interface should honk")

	// Hopefully make this one work one day. I think it might be impossible?
	//var instance3 *CarService
	//app.Resolve(instance3)
	//assert.Equal(t, "Env Honk!", instance3.Honk(), "Zero-value pointers should work")
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
