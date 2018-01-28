package templates

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Renderer struct {
	root     string
	funcs    template.FuncMap
	views    map[string]map[string]*template.Template
	partials map[string]*template.Template
}

func NewRenderer(root string, funcs template.FuncMap) *Renderer {
	return &Renderer{root: root, funcs: funcs}
}

func (r *Renderer) Init() error {
	templatePaths, err := r.findTemplates(r.root, "views")
	if err != nil {
		return err
	}

	layoutPaths, err := r.findTemplates(r.root, "layouts")
	if err != nil {
		return err
	}

	// Create all possible combinations of views to bases
	r.views = map[string]map[string]*template.Template{}
	for _, templatePath := range templatePaths {
		viewName := strings.TrimPrefix(templatePath, filepath.Join(r.root, "views")+"/")
		r.views[viewName] = map[string]*template.Template{}
		for _, layoutPath := range append(layoutPaths, "") {
			var (
				layoutName string
				t          *template.Template
				err        error
			)

			if layoutPath == "" {
				layoutName = ""
				t, err = template.New("").Funcs(r.funcs).ParseFiles(templatePath)
			} else {
				// NOTE: Layout must be parsed before template so {{ block }} defaults work
				layoutName = strings.TrimPrefix(layoutPath, filepath.Join(r.root, "layouts")+"/")
				t, err = template.
					New(viewName+"::"+layoutName).
					Funcs(r.funcs).
					ParseFiles(layoutPath, templatePath)
			}

			if err != nil {
				return err
			}

			r.views[viewName][layoutName] = t
		}
	}

	partialPaths, err := r.findTemplates(r.root, "partials")
	if err != nil {
		return err
	}

	r.partials = make(map[string]*template.Template)
	for _, partialPath := range partialPaths {
		base := filepath.Join(r.root, "partials")
		name := strings.TrimPrefix(partialPath, base+"/")
		tmpl, err := ioutil.ReadFile(partialPath)
		if err != nil {
			return err
		}
		t, err := template.New(name).Funcs(r.funcs).Parse(string(tmpl))
		if err != nil {
			return err
		}
		r.partials[name] = t
	}

	fmt.Printf(
		"[renderer] Parsed %d views and %d partials\n",
		len(templatePaths),
		len(partialPaths),
	)

	return nil
}

func (r *Renderer) RenderString(w io.Writer, html string, data interface{}) error {
	t, err := template.New("").Funcs(r.funcs).Parse(html)
	if err != nil {
		return err
	}

	return t.Execute(w, data)
}

func (r *Renderer) RenderPartial(name string, data interface{}) string {
	t, ok := r.partials[name]
	if !ok {
		panic("Partial not found with name " + name)
	}

	var w bytes.Buffer
	err := t.Execute(&w, data)
	if err != nil {
		panic("Failed to render partial: " + err.Error())
	}

	return w.String()
}

func (r *Renderer) RenderTemplate(w io.Writer, template, layout string, data interface{}) error {
	if len(r.views) == 0 {
		return errors.New(fmt.Sprintf("No views found in %s", r.root))
	}

	if _, ok := r.views[template]; !ok {
		templates := make([]string, 0)
		for name := range r.views {
			templates = append(templates, name)
		}
		options := strings.Join(templates, ", ")
		return errors.New(fmt.Sprintf("Template not found '%s'. Options are %s", template, options))
	}

	t, ok := r.views[template][layout]
	if !ok {
		return errors.New(fmt.Sprintf("Layout (%s) not found. Options %#v", layout, r.views[template]))
	}

	if layout == "" {
		return t.ExecuteTemplate(w, template, data)
	} else {
		return t.ExecuteTemplate(w, layout, data)
	}
}

func (r *Renderer) findTemplates(dirs ...string) ([]string, error) {
	dir := filepath.Join(dirs...)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return make([]string, 0), nil
	}

	fileInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0)
	for _, f := range fileInfo {
		p := filepath.Join(dir, f.Name())
		if f.IsDir() {
			// Recurse if directory
			more, err := r.findTemplates(p)
			if err != nil {
				return nil, err
			}
			paths = append(paths, more...)
		} else {
			paths = append(paths, p)
		}
	}

	return paths, nil
}
