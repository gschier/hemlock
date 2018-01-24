package funcs

import (
	"github.com/gschier/hemlock"
	"github.com/gschier/hemlock/interfaces"
	"html/template"
	"log"
	"net/url"
	"path"
)

func fnURL (name string) template.URL {
	base := hemlock.App().Config.URL

	u, err := url.Parse(base)
	if err != nil {
		log.Panicf("Invalid App URL: %s", base)
	}

	var router interfaces.Router
	hemlock.App().Resolve(&router)

	u.Path = path.Join(u.Path, router.URL(name, nil))
	return template.URL(u.String())
}
