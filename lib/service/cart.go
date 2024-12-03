package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"vshop/lib/model"

	"github.com/alexedwards/scs/v2"
)

type CartItem struct {
	IProdID int32
	IQty    int32
}

type Cart []CartItem

func (s *Service) CartCount(r *http.Request) (int, error) {

	var cart = []CartItem{}
	var ok bool
	if s.SessionManager.Exists(r.Context(), "cart") {
		c := s.SessionManager.Get(r.Context(), "cart")
		cart, ok = c.([]CartItem)
		if !ok {
			return 0, errors.New("cannot cast session cart to type cart")
		}
	}
	cartCount := len(cart)

	return cartCount, nil
}

func (s *Service) CartProducts(r *http.Request) ([]model.CartProduct, error) {

	if !s.SessionManager.Exists(r.Context(), "cart") {
		return []model.CartProduct{}, nil
	}

	var cart = []CartItem{}
	var ok bool
	c := s.SessionManager.Get(r.Context(), "cart")
	cart, ok = c.([]CartItem)
	if !ok {
		return nil, errors.New("cannot cast session cart to type cart")
	}

	prodIDs := []int32{}
	id2qty := make(map[int32]int32)

	for _, item := range cart {
		id2qty[item.IProdID] = item.IQty
		prodIDs = append(prodIDs, item.IProdID)
	}

	var err error
	products, err := s.Model.CartProducts(prodIDs)
	if err != nil {
		return nil, fmt.Errorf("error fetching products from cart: %w", err)
	}

	cartProducts := make([]model.CartProduct, 0)
	for i, product := range products {
		cartProducts = append(cartProducts, model.CartProduct{
			Product:       product,
			ICartQuantity: id2qty[product.IProdID],
			FAmount:       product.FPrice * float64(id2qty[product.IProdID]),
		})
		cartProducts[i].ICartQuantity = id2qty[product.IProdID]
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

type CartRequest struct {
	ID  string `json:"iProdID"`
	Qty string `json:"iQty"`
}

const MAX_REQUEST_BODY_BYTES = 100 * 1024

func (s *Service) HandleAddToCart(w http.ResponseWriter, r *http.Request) error {

	cartReq := new(CartRequest)
	reader := http.MaxBytesReader(w, r.Body, MAX_REQUEST_BODY_BYTES)
	err := json.NewDecoder(reader).Decode(cartReq)
	if err != nil {
		return fmt.Errorf("error decoding request body: %w", err)
	}

	iProdID, err := strconv.Atoi(cartReq.ID)
	if err != nil {
		return fmt.Errorf("error converting id %s to int: %w", cartReq.ID, err)
	}

	iQty, err := strconv.Atoi(cartReq.Qty)
	if err != nil {
		return fmt.Errorf("error converting quantity %s to int: %w", cartReq.Qty, err)
	}

	if err := s.AddToCart(int32(iProdID), int32(iQty), r); err != nil {
		return fmt.Errorf("error adding to cart: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "plain/text")
	w.Write([]byte("worked"))

	return nil
}
