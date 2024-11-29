package service

import (
	"net/http"
	"vshop/lib/model"
)

type IndexPageData struct {
	TopCategories []*model.Category
	PathToCurrent []*model.Category
}

func (s *Service) Index(w http.ResponseWriter, r *http.Request) error {

	catRoot, err := s.Model.CategoryTree()
	if err != nil {
		return err
	}
	pathToCurrent := []*model.Category{catRoot}
	topCategories := catRoot.Children

	data := IndexPageData{
		TopCategories: topCategories,
		PathToCurrent: pathToCurrent,
	}

	if err := s.render(w, "index.go.html", data); err != nil {
		return err
	}

	return nil
}
