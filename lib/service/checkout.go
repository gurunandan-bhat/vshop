package service

import (
	"fmt"
	"net/http"
	"vshop/lib/model"
)

func (s *Service) Checkout(w http.ResponseWriter, r *http.Request) error {

	cartCount, err := s.CartCount(r)
	if err != nil {
		return fmt.Errorf("error fetching cart products: %w", err)
	}

	cartProducts, err := s.CartProducts(r)
	if err != nil {
		return err
	}

	catRoot, err := s.Model.CategoryTree()
	if err != nil {
		return err
	}

	data := struct {
		TopCategories []*model.Category
		CartProducts  []model.Product
		CartCount     int
		IncludeJS     bool
	}{
		TopCategories: catRoot.Children,
		CartProducts:  cartProducts,
		CartCount:     cartCount,
		IncludeJS:     false,
	}

	if err := s.render(w, "checkout.go.html", data); err != nil {
		return err
	}

	return nil
}
