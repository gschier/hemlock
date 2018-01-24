package funcs

import "html/template"

var Funcs = template.FuncMap{
	"asset": fnAsset,
	"url": fnURL,
	"partial": fnPartial,
}
