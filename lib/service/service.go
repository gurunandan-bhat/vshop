package service

import (
	"html/template"
	"log"
	"net/http"
	"vshop/lib/config"
	"vshop/lib/cstore"
	"vshop/lib/model"
	"vshop/lib/render"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/csrf"
)

type Service struct {
	Model         *model.Model
	Muxer         *chi.Mux
	SessionStore  *cstore.SessionStore
	StaticDir     string
	S3Root        string
	TemplateCache map[string]*template.Template
}

func NewService(cfg *config.Config) (*Service, error) {

	mux := chi.NewRouter()

	// force a redirect to https:// in production
	if cfg.InProduction {
		mux.Use(middleware.SetHeader(
			"Strict-Transport-Security",
			"max-age=63072000; includeSubDomains",
		))
	}

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	csrfMiddleware := csrf.Protect(
		[]byte(cfg.Security.CSRFKey),
		csrf.Secure(cfg.InProduction),
		csrf.SameSite(csrf.SameSiteStrictMode),
	)
	mux.Use(csrfMiddleware)

	store := cstore.NewStore(cfg)

	model, err := model.NewModel(cfg)
	if err != nil {
		log.Fatalf("Error initializing database connection: %s", err)
	}

	tmplCache, err := render.NewTemplates(cfg.TemplateRoot)
	if err != nil {
		log.Fatalf("Cannot build template cache: %s", err)
	}

	s := &Service{
		SessionStore:  store,
		Model:         model,
		Muxer:         mux,
		StaticDir:     "./static",
		S3Root:        cfg.S3Root,
		TemplateCache: tmplCache,
	}

	s.setRoutes()

	return s, nil
}

func (s *Service) setRoutes() {

	s.Muxer.Method(http.MethodGet, "/", ServiceHandler(s.Index))
	s.Muxer.Method(http.MethodGet, "/static/*", ServiceHandler(s.Static))

	s.Muxer.Method(http.MethodGet, "/category/{vUrlName}", ServiceHandler(s.CategoryProducts))
	s.Muxer.Method(http.MethodGet, "/product/{vUrlName}", ServiceHandler(s.Product))
	// s.Muxer.Method(http.MethodGet, "/register", ServiceHandler(s.Register))
	// s.Muxer.Method(http.MethodGet, "/tree/{vUrlName}", ServiceHandler(s.Tree))
	// s.Muxer.Method(http.MethodGet, "/quick-view/product/{iProdID}", ServiceHandler(s.ProductImages))
}
