package providers

import "github.com/gschier/hemlock"

type RouteProvider struct {

}

func (*RouteProvider) Register(*hemlock.Container) {
	panic("implement me")
}

func (*RouteProvider) Boot(*hemlock.Application) {
	panic("implement me")
}


