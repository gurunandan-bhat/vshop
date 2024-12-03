package service

import (
	"fmt"
	"net/http"
	"vshop/lib/config"
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

	totalAmt := 0
	for _, p := range cartProducts {
		totalAmt += int(p.FAmount)
	}

	catRoot, err := s.Model.CategoryTree()
	if err != nil {
		return err
	}

	cfg, err := config.Configuration()
	if err != nil {
		return err
	}

	data := struct {
		TopCategories []*model.Category
		CartProducts  []model.CartProduct
		CartCount     int
		TotalAmount   int
		IncludeJS     bool
		CheckoutURL   string
	}{
		TopCategories: catRoot.Children,
		CartProducts:  cartProducts,
		CartCount:     cartCount,
		TotalAmount:   totalAmt,
		IncludeJS:     false,
		CheckoutURL:   cfg.CheckoutURL,
	}

	if err := s.render(w, "checkout.go.html", data); err != nil {
		return err
	}

	return nil
}
