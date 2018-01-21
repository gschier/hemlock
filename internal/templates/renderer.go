package templates

import (
	"fmt"
	"github.com/pkg/errors"
	"html/template"
	"io"
	"io/ioutil"
	"path/filepath"
)

type Renderer struct {
	root      string
	templates map[string]map[string]*template.Template
}

func NewRenderer(root string) *Renderer {
	r := &Renderer{root: root}
	r.Load()
	return r
}

func (r *Renderer) Load() error {
	templatePaths, err := r.findTemplates(r.root)
	if err != nil {
		return err
	}

	layoutPaths, err := r.findTemplates(r.root, "layouts")
	if err != nil {
		return err
	}

	partialPaths, err := r.findTemplates(r.root, "partials")
	if err != nil {
		return err
	}

	// Create all possible combinations of templates to bases
	r.templates = map[string]map[string]*template.Template{}
	for _, templatePath := range templatePaths {
		templateName := filepath.Base(templatePath)
		r.templates[templateName] = map[string]*template.Template{}
		for _, layoutPath := range layoutPaths {
			paths := append(partialPaths, templatePath, layoutPath)
			t, err := template.ParseFiles(paths...)
			if err != nil {
				return err
			}

			layoutName := filepath.Base(layoutPath)
			r.templates[templateName][layoutName] = t
		}
	}

	fmt.Printf("[renderer] Parsed %d templates\n", len(templatePaths))

	return nil
}

func (r *Renderer) RenderTemplate(w io.Writer, name, layout string, data interface{}) error {
	templateName := name + ".html"
	layoutName := layout + ".html"
	if _, ok := r.templates[templateName]; !ok {
		return errors.New("Template not found: " + name)
	}

	t, ok := r.templates[templateName][layoutName]
	if !ok {
		return errors.New("Layout not found: " + layout)
	}

	return t.ExecuteTemplate(w, layoutName, data)
}

func (r *Renderer) findTemplates(dirs ...string) ([]string, error) {
	dir := filepath.Join(dirs...)
	fileInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0)
	for _, f := range fileInfo {
		if f.IsDir() {
			// TODO: Implement recursive search
			continue
		}
		paths = append(paths, filepath.Join(dir, f.Name()))
	}

	return paths, nil
}
