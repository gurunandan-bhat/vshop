package service

import (
	"fmt"
	"net/http"
	"vshop/lib/model"

	"github.com/go-chi/chi/v5"
)

type CategoryPageData struct {
	TopCategories    []*model.Category
	CurrentCategory  *model.Category
	PathToCurrent    []*model.Category
	CategorySideBar  []*model.Category
	CategoryProducts []model.CategoryProducts
	S3Root           string
}

func (s *Service) CategoryProducts(w http.ResponseWriter, r *http.Request) error {

	vUrlName := chi.URLParam(r, "vUrlName")

	catRoot, err := s.Model.CategoryTree()
	if err != nil {
		return err
	}
	topCategories := catRoot.Children

	catMap, err := s.Model.CategoryMapByURL()
	if err != nil {
		return err
	}

	currentCategory, found := catMap[vUrlName]
	if !found {
		return fmt.Errorf("error fetching (sub)-tree root with URL: %s", vUrlName)
	}
	if currentCategory.IParentID == 0 && len(currentCategory.Children) > 0 {
		for {
			c := currentCategory.Children[0]
			currentCategory = c
			if currentCategory.IProductCount > 0 {
				break
			}
		}
	}

	pathToCurrent, err := s.Model.ReverseCatTree(currentCategory)
	if err != nil {
		return err
	}

	top := pathToCurrent[0]
	sideBar, err := s.Model.AllSubCategories(top, true)
	if err != nil {
		return err
	}

	products, err := s.Model.CategoryProducts(vUrlName)
	if err != nil {
		return err
	}

	data := CategoryPageData{
		TopCategories:    topCategories,
		CurrentCategory:  currentCategory,
		PathToCurrent:    pathToCurrent,
		CategorySideBar:  sideBar,
		CategoryProducts: *products,
		S3Root:           s.S3Root,
	}

	if err := s.render(w, "category.go.html", data); err != nil {
		return err
	}

	return nil
}
