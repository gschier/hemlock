package funcs

import (
	"github.com/gschier/hemlock"
	"html/template"
	"log"
	url2 "net/url"
	"path"
	"strings"
)

func asset(app *hemlock.Application) interface{} {
	return func(name string) template.URL {
		base := app.Config.URL
		publicDir := app.Config.PublicDirectory
		publicPrefix := app.Config.PublicPrefix
		fullURL := name
		if strings.Contains(base, "://") {
			u, err := url2.Parse(app.Config.URL)
			if err != nil {
				log.Panicf("Invalid App URL: %s", base)
			}
			u.Path = path.Join(u.Path, publicPrefix, publicDir, name)
			fullURL = u.String()
		} else {
			fullURL = path.Join(base, publicPrefix, publicDir, name)
		}

		return template.URL(fullURL)
	}
}
