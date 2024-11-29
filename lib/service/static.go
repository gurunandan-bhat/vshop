package service

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (s *Service) Static(w http.ResponseWriter, r *http.Request) error {

	rctx := chi.RouteContext(r.Context())
	pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")

	root := http.Dir(s.StaticDir)
	fs := http.StripPrefix(pathPrefix, http.FileServer(root))

	fs.ServeHTTP(w, r)

	return nil
}
