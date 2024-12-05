package service

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"vshop/lib/model"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
)

func (s *Service) Product(w http.ResponseWriter, r *http.Request) error {

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			return fmt.Errorf("error parsing form: %w", err)
		}
		values := r.PostForm
		p, ok := values["iProdID"]
		if !ok {
			return fmt.Errorf("product is not available for update: %#v", values)
		}
		iProdID, err := strconv.Atoi(p[0])
		if err != nil {
			return fmt.Errorf("product has the wrong type: %s: %w", p[0], err)
		}

		q, ok := values["iQty"]
		if !ok {
			return fmt.Errorf("product quantity is not available for update: %#v", values)
		}
		iQty, err := strconv.Atoi(q[0])
		if err != nil {
			return fmt.Errorf("product quantity has the wrong type: %s: %w", q[0], err)
		}

		if !((iProdID > 0) && (iQty > 0)) {
			return errors.New("invalid cart item, id and/or qty cannot be zero")
		}

		var cart = []CartItem{}
		if s.SessionManager.Exists(r.Context(), "cart") {
			c := s.SessionManager.Get(r.Context(), "cart")
			cart, ok = c.([]CartItem)
			if !ok {
				return errors.New("cannot cast session cart to type cart")
			}
		}

		cart = append(cart, CartItem{int32(iProdID), int32(iQty)})
		s.SessionManager.Put(r.Context(), "cart", cart)

		// check that the session was indeed modified
		if s.SessionManager.Status(r.Context()) != scs.Modified {
			return fmt.Errorf("session was not modified after adding %d of product id %d", iQty, iProdID)
		}

		product, err := s.Model.Product(int32(iProdID))
		if err != nil {
			return fmt.Errorf("error fetching product by ID %d: %w", iProdID, err)
		}

		http.Redirect(w, r, fmt.Sprintf("/product/%s", product.VURLName), http.StatusMovedPermanently)
	}

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

	cartCount, err := s.CartCount(r)
	if err != nil {
		return fmt.Errorf("error fetching cart products: %w", err)
	}

	data := struct {
		model.Product
		Attributes    []model.ProductAttribute
		TopCategories []*model.Category
		PathToCurrent any
		S3Root        string
		CartCount     int
		IncludeJS     bool
	}{
		Product:       *product,
		Attributes:    *attribs,
		TopCategories: catRoot.Children,
		PathToCurrent: []*model.Category{catRoot},
		S3Root:        s.S3Root,
		CartCount:     cartCount,
		IncludeJS:     false,
	}

	if err := s.render(w, "product.go.html", data); err != nil {
		return err
	}

	return nil
}
