package funcs

import (
	"github.com/gschier/hemlock"
	"html/template"
)

func Funcs (app *hemlock.Application) *template.FuncMap {
	return &template.FuncMap{
		"asset": asset(app),
		"url": url(app),
		"partial": partial(app),
		"route": route(app),
	}
}

