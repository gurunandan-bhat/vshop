package service

import (
	"fmt"
	"net/http"
)

func (s *Service) render(w http.ResponseWriter, page string, data any) error {

	// Retrieve the appropriate template set from the cache based on the page
	// name (like 'home.tmpl'). If no entry exists in the cache with the
	// provided name, then create a new error and call the serverError() helper
	// method that we made earlier and return.
	ts, ok := s.TemplateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		return err
	}

	// Write out the provided HTTP status code ('200 OK', '400 Bad Request' etc).
	w.WriteHeader(http.StatusOK)
	// Execute the template set and write the response body. Again, if there
	// is any error we call the serverError() helper.
	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	return nil
}
