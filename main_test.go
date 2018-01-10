package begonia_test

import (
	"github.com/gschier/begonia"
	"fmt"
)

func NewTestApplication(p ...begonia.Provider) *begonia.Application {
	return begonia.NewApplication(&begonia.Config{
		ApplicationConfig: &begonia.ApplicationConfig{Providers: p},
	})
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
