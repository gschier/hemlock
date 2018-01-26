package funcs

import (
	"bytes"
	"github.com/gschier/hemlock"
	"github.com/gschier/hemlock/internal/templates"
	"html/template"
	"io/ioutil"
)

func partial (name string, data ...interface{}) template.HTML {
	app := hemlock.App()
	var renderer templates.Renderer
	app.Resolve(&renderer)

	partialPath := app.Path(app.Config.TemplatesDirectory, "partials", name)
	content, err := ioutil.ReadFile(partialPath)
	if err != nil {
		panic(err)
	}

	var d interface{} = nil
	if len(data) > 0 {
		d = data[0]
	}

	var w bytes.Buffer
	err = renderer.RenderString(&w, string(content), d)
	if err != nil {
		panic(err)
	}

	return template.HTML(w.String())
}
