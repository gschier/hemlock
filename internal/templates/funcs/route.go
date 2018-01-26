package funcs

import (
	"github.com/gschier/hemlock"
	"github.com/gschier/hemlock/interfaces"
	"html/template"
	"strings"
)

func route(name string, params ...string) template.URL {
	var router interfaces.Router
	hemlock.App().Resolve(&router)

	// Split name=value pairs into map
	paramsMap := make(map[string]string)
	for _, p := range params {
		v := strings.SplitN(p, "=", 2)
		// TODO: Better error messages here
		paramsMap[v[0]] = v[1]
	}

	return template.URL(router.Route(name, paramsMap))
}
