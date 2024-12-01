package service

import (
	"fmt"
	"net/http"
	"vshop/lib/model"
)

type IndexPageData struct {
	TopCategories []*model.Category
	PathToCurrent []*model.Category
	CartCount     int
	IncludeJS     bool
}

func (s *Service) Index(w http.ResponseWriter, r *http.Request) error {

	catRoot, err := s.Model.CategoryTree()
	if err != nil {
		return err
	}
	pathToCurrent := []*model.Category{catRoot}
	topCategories := catRoot.Children

	cartCount, err := s.CartCount(r)
	if err != nil {
		return fmt.Errorf("error fetching cart products: %w", err)
	}

	data := IndexPageData{
		TopCategories: topCategories,
		PathToCurrent: pathToCurrent,
		CartCount:     cartCount,
		IncludeJS:     false,
	}

	if err := s.render(w, "index.go.html", data); err != nil {
		return err
	}

	return nil
}
