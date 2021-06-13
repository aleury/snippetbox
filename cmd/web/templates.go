package main

import (
	"html/template"
	"path/filepath"
	"time"

	"adameury.io/snippetbox/pkg/models"
)

type defaultData struct {
	CurrentYear int
}

type homeTemplateData struct {
	defaultData
	Snippets []models.Snippet
}

type showTemplateData struct {
	defaultData
	Snippet models.Snippet
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var tmplFuncs = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// Add page to a template set
		ts, err := template.New(name).Funcs(tmplFuncs).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add layouts to the template set
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// Add partials to the template set
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
