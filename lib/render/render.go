package render

import (
	"fmt"
	"html/template"
	"path/filepath"
)

var mathSub = func(i, d int) int {

	return i - d
}

var tmplFuncs = template.FuncMap{
	"sub": mathSub,
}

func NewTemplates(templateRoot string) (map[string]*template.Template, error) {

	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	// Use the filepath.Glob() function to get a slice of all filepaths that
	// match the pattern "./ui/html/pages/*.tmpl". This will essentially gives
	// us a slice of all the filepaths for our application 'page' templates
	// like: [ui/html/pages/home.tmpl ui/html/pages/view.tmpl]
	pages, err := filepath.Glob(fmt.Sprintf("%s/pages/*.go.html", templateRoot))
	if err != nil {
		return nil, err
	}

	partials, err := filepath.Glob(fmt.Sprintf("%s/partials/*.go.html", templateRoot))
	if err != nil {
		return nil, err
	}

	base := fmt.Sprintf("%s/base.go.html", templateRoot)

	var templates []string
	templates = append(append(templates, base), partials...)

	// Loop through the page filepaths one-by-one.
	for _, page := range pages {

		// Extract the file name (like 'home.tmpl') from the full filepath
		// and assign it to the name variable.
		name := filepath.Base(page)

		// Create a slice containing the filepaths for our base template, any
		// partials and the page.
		files := append(templates, page)

		// Parse the files into a template set.
		ts, err := template.New("base").Funcs(tmplFuncs).ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map, using the name of the page
		// (like 'home.tmpl') as the key.
		cache[name] = ts
	}
	// Return the map.
	return cache, nil
}
