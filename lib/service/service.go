package service

import (
	"encoding/gob"
	"html/template"
	"log"
	"net/http"
	"vshop/lib/config"
	"vshop/lib/model"
	"vshop/lib/render"
	"vshop/lib/scsstore.go"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Service struct {
	Model          *model.Model
	Muxer          *chi.Mux
	SessionManager *scs.SessionManager
	StaticDir      string
	S3Root         string
	TemplateCache  map[string]*template.Template
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

	model, err := model.NewModel(cfg)
	if err != nil {
		log.Fatalf("Error initializing database connection: %s", err)
	}
	gob.Register([]CartItem{})
	scsstore, err := scsstore.NewSCSStore(model)
	if err != nil {
		log.Fatalf("Error initializing session store: %s", err)
	}

	mux.Use()
	tmplCache, err := render.NewTemplates(cfg.TemplateRoot)
	if err != nil {
		log.Fatalf("Cannot build template cache: %s", err)
	}

	s := &Service{
		SessionManager: scsstore,
		Model:          model,
		Muxer:          mux,
		StaticDir:      "./static",
		S3Root:         cfg.S3Root,
		TemplateCache:  tmplCache,
	}

	s.setRoutes()

	return s, nil
}

func (s *Service) setRoutes() {

	s.Muxer.Method(http.MethodGet, "/static/*", ServiceHandler(s.Static))

	s.Muxer.Group(func(r chi.Router) {

		r.Use(s.SessionManager.LoadAndSave)

		r.Method(http.MethodGet, "/", ServiceHandler(s.Index))
		r.Method(http.MethodGet, "/category/{vUrlName}", ServiceHandler(s.CategoryProducts))
		r.Method(http.MethodGet, "/product/{vUrlName}", ServiceHandler(s.Product))
		r.Method(http.MethodGet, "/checkout", ServiceHandler(s.Checkout))

		r.Method(http.MethodPost, "/add-to-cart", ServiceHandler(s.HandleAddToCart))

		// r.Method(http.MethodGet, "/register", ServiceHandler(s.Register))
		// r.Method(http.MethodGet, "/tree/{vUrlName}", ServiceHandler(s.Tree))
		// r.Method(http.MethodGet, "/quick-view/product/{iProdID}", ServiceHandler(s.ProductImages))
	})
}
