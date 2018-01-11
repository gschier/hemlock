package testutil

import (
	"github.com/gschier/hemlock"
	"fmt"
)

func NewTestApplication(p ...hemlock.Provider) *hemlock.Application {
	return hemlock.NewApplication(&hemlock.Config{
		ApplicationConfig: &hemlock.ApplicationConfig{Providers: p},
	})
}

type CarServiceInterface interface {
	Honk() string
}

type CarService struct {
	Noise string
}

func (cs *CarService) Honk() string {
	if cs.Noise == "" {
		return "Honk!"
	}

	return cs.Noise
}

type CarServiceProvider struct {
	noise string
}

func (csp *CarServiceProvider) Register(ioc *hemlock.Container) {
	ioc.Singleton(func(app *hemlock.Application) (*CarService, error) {
		return &CarService{Noise: app.Env("honk")}, nil
	})
}

func (csp *CarServiceProvider) Boot(app *hemlock.Application) {
	carService := app.Make(new(CarServiceInterface)).(*CarService)
	fmt.Printf("HELLO? %v\n", carService.Honk())
}
