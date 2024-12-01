package service

import (
	"errors"
	"fmt"
	"net/http"
	"vshop/lib/model"

	"github.com/alexedwards/scs/v2"
)

type CartItem struct {
	IProdID int32
	IQty    int32
}

type Cart []CartItem

func (s *Service) GetCartProducts(r *http.Request) ([]model.Product, error) {

	var cartProducts []model.Product

	var cart = []CartItem{}
	var ok bool
	if s.SessionManager.Exists(r.Context(), "cart") {
		c := s.SessionManager.Get(r.Context(), "cart")
		cart, ok = c.([]CartItem)
		if !ok {
			return nil, errors.New("cannot cast session cart to type cart")
		}

		prodIDs := []int32{}
		for _, item := range cart {
			prodIDs = append(prodIDs, item.IProdID)
		}

		var err error
		cartProducts, err = s.Model.GetCartProducts(prodIDs)
		if err != nil {
			return nil, fmt.Errorf("error fetching products from cart: %w", err)
		}
	}
	return cartProducts, nil
}

func (s *Service) AddToCart(iProdID, iQty int32, r *http.Request) error {

	if !((iProdID > 0) && (iQty > 0)) {
		return errors.New("invalid cart item, id and/or qty cannot be zero")
	}

	var cart = []CartItem{}
	var ok bool
	if s.SessionManager.Exists(r.Context(), "cart") {
		c := s.SessionManager.Get(r.Context(), "cart")
		cart, ok = c.([]CartItem)
		if !ok {
			return errors.New("cannot cast session cart to type cart")
		}
	}

	cart = append(cart, CartItem{iProdID, iQty})
	s.SessionManager.Put(r.Context(), "cart", cart)

	// check that the session was indeed modified
	if s.SessionManager.Status(r.Context()) != scs.Modified {
		return fmt.Errorf("session was not modified after adding %d of product id %d", iQty, iProdID)
	}

	return nil
}
