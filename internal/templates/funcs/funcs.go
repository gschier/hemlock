package funcs

import "html/template"

var Funcs = template.FuncMap{
	"asset": asset,
	"url": url,
	"partial": partial,
	"route": route,
}
