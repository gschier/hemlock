package funcs

import (
	"github.com/gschier/hemlock"
	"html/template"
	"log"
	"net/url"
	"path"
	"strings"
)

func fnAsset(name string) template.URL {
	app := hemlock.App()
	base := app.Config.URL
	assetBase := app.Config.AssetBase
	fullURL := name
	if strings.Contains(base, "://") {
		u, err := url.Parse(app.Config.URL)
		if err != nil {
			log.Panicf("Invalid App URL: %s", base)
		}
		u.Path = path.Join(u.Path, assetBase, name)
		fullURL = u.String()
	} else {
		fullURL = path.Join(base, assetBase, name)
	}

	return template.URL(fullURL)
}
