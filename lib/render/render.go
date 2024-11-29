package render

import (
	"fmt"
	"html/template"
	"net/http"
)

type Templates struct {
	Root  string
	cache *template.Template
}

func NewTemplates(templateRoot string) (*Templates, error) {

	tc, err := template.ParseGlob(templateRoot + "/*.go.html")
	if err != nil {
		return nil, err
	}

	return &Templates{
		Root:  templateRoot,
		cache: tc,
	}, nil
}

func (t *Templates) Render(w http.ResponseWriter, template string, data any) error {

	// Check whether that template exists in the cache
	if t.cache.Lookup(template) == nil {
		return fmt.Errorf("template %s is not available in the cache", template)
	}

	if err := t.cache.ExecuteTemplate(w, template, data); err != nil {
		fmt.Println(template, err)
		return err
	}

	return nil
}
