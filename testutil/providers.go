package testutil

import (
	"github.com/gschier/hemlock"
	"github.com/gschier/hemlock/interfaces"
)

func NewTestApplication(p ...hemlock.Provider) *hemlock.Application {
	return hemlock.NewApplication(&hemlock.Config{}, p)
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

func (csp *CarServiceProvider) Register(ioc interfaces.Container) {
	ioc.Singleton(func(app *hemlock.Application) (*CarService, error) {
		return &CarService{Noise: app.Env("honk")}, nil
	})
}

func (csp *CarServiceProvider) Boot(app *hemlock.Application) error {
	return nil
}

type StringServiceProvider struct{}

func (ssp *StringServiceProvider) Register(ioc interfaces.Container) {
	ioc.Singleton(func(app *hemlock.Application) (string, error) {
		return "Hello!", nil
	})
}

func (ssp *StringServiceProvider) Boot(app *hemlock.Application) error {
	return nil
}
