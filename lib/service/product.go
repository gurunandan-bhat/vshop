package service

import (
	"errors"
	"fmt"
	"net/http"
	"vshop/lib/model"

	"github.com/go-chi/chi/v5"
)

func (s *Service) Product(w http.ResponseWriter, r *http.Request) error {

	vUrlName := chi.URLParam(r, "vUrlName")
	if vUrlName == "" {
		return errors.New("Invalid Product Name")
	}

	product, err := s.Model.ProductByUrl(vUrlName)
	if err != nil {
		return fmt.Errorf("error fetching product: %w", err)
	}

	attribs, err := s.Model.ProductAttributes(product.IProdID)
	if err != nil {
		return fmt.Errorf("error fetching product attributes: %w", err)
	}

	catRoot, err := s.Model.CategoryTree()
	if err != nil {
		return err
	}

	cartProducts, err := s.GetCartProducts(r)
	if err != nil {
		return fmt.Errorf("error fetching cart products: %w", err)
	}

	data := struct {
		model.Product
		Attributes    []model.ProductAttribute
		TopCategories []*model.Category
		PathToCurrent any
		S3Root        string
		CartProducts  []model.Product
	}{
		Product:       *product,
		Attributes:    *attribs,
		TopCategories: catRoot.Children,
		PathToCurrent: []*model.Category{catRoot},
		S3Root:        s.S3Root,
		CartProducts:  cartProducts,
	}

	if err := s.render(w, "product.go.html", data); err != nil {
		return err
	}

	return nil
}
