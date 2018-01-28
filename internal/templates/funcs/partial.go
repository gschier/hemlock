package funcs

import (
	"github.com/gschier/hemlock"
	"github.com/gschier/hemlock/internal/templates"
	"html/template"
)

func partial(app *hemlock.Application) interface{} {
	return func(name string, data ...interface{}) template.HTML {
		var renderer templates.renderer
		app.Resolve(&renderer)

		var renderData interface{}
		if len(data) > 0 {
			renderData = data[0]
		}

		return template.HTML(renderer.RenderPartial(name, renderData))
	}
}
