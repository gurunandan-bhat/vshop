package service

import (
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
)

func (s *Service) Register(w http.ResponseWriter, r *http.Request) error {

	tmpl, err := template.New("register").ParseFiles(
		"./templates/header_offer.go.html",
		"./templates/register.go.html",
	)
	if err != nil {
		return err
	}

	data := map[string]any{
		csrf.TemplateTag: csrf.TemplateField(r),
	}
	tmpl.Execute(w, data)

	return nil
}

type Registrant struct {
	Email        string `schema:"email, required"`
	Password     string `schema:"password, required"`
	PasswordCopy string `schema:"password-copy, required"`
	Subscribe    bool   `schema:"subscribe"`
}

func (s *Service) RegisterNewUser(w http.ResponseWriter, r *http.Request) error {

	return nil
}
