package funcs

import (
	"github.com/gschier/hemlock"
	"github.com/gschier/hemlock/interfaces"
	"html/template"
)

func url(path string) template.URL {
	var router interfaces.Router
	hemlock.App().Resolve(&router)
	return template.URL(router.URL(path))
}
