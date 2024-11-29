package service

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Service) Tree(w http.ResponseWriter, r *http.Request) error {

	vUrlName := chi.URLParam(r, "vUrlName")
	catMap, err := s.Model.CategoryMapByURL()
	if err != nil {
		return err
	}

	start, found := catMap[vUrlName]
	if !found {
		return fmt.Errorf("error fetching (sub)-tree root with URL: %s", vUrlName)
	}

	s.Template.Render(w, "tree", start)

	return nil
}
