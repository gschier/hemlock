package funcs

import (
	"github.com/gschier/hemlock"
	"html/template"
	url2 "net/url"
	"path"
	"strings"
)

func asset(app *hemlock.Application) interface{} {
	return func(name string) template.URL {
		var (
			err     error
			fullURL *url2.URL
		)

		publicPrefix := app.Config.PublicPrefix

		prefixAbsolute := strings.HasPrefix(publicPrefix, "https://") ||
			strings.HasPrefix(publicPrefix, "http://") ||
			strings.HasPrefix(publicPrefix, "//")

		if prefixAbsolute {
			fullURL, err = url2.Parse(publicPrefix)
			if err != nil {
				panic(err)
			}
		} else {
			fullURL, err = url2.Parse(app.Config.URL)
			if err != nil {
				panic(err)
			}
			fullURL.Path = path.Join(fullURL.Path, publicPrefix)
		}

		fullURL.Path = path.Join(fullURL.Path, name)

		if app.IsProd() {
			q := fullURL.Query()
			q.Set("v", hemlock.CacheBustKey)
			fullURL.RawQuery = q.Encode()
		}

		return template.URL(fullURL.String())
	}
}
