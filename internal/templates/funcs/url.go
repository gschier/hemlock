package funcs

import (
	"github.com/gschier/hemlock"
	"github.com/gschier/hemlock/interfaces"
	"html/template"
)

func url(app *hemlock.Application) interface{} {
	return func(path string) template.URL {
		var router interfaces.Router
		app.Resolve(&router)
		return template.URL(router.URL(path))
	}
}
