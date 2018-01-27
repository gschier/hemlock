package funcs

import (
	"github.com/gschier/hemlock"
	"github.com/gschier/hemlock/internal/templates"
	"html/template"
)

func partial (name string, data ...interface{}) template.HTML {
	app := hemlock.App()
	var renderer templates.Renderer
	app.Resolve(&renderer)

	var renderData interface{}
	if len(data) > 0 {
		renderData = data[0]
	}

	return template.HTML(renderer.RenderPartial(name, renderData))
}
