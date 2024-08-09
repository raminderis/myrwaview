package views

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/raminderis/myrwaview/migrations/context"
	"github.com/raminderis/myrwaview/models"
)

type Template struct {
	htmlTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}) {
	htmlTpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("cloning template: %v", err)
		http.Error(w, "There was error rendering the page", http.StatusInternalServerError)
		return
	}
	htmlTpl = htmlTpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
			"currentUser": func() *models.User {
				return context.User(r.Context())
			},
		},
	)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var buf bytes.Buffer
	err = htmlTpl.Execute(&buf, data)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was error executing the template", http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

func Parse(filepath string) (Template, error) {
	htmlTpl, err := template.ParseFiles(filepath)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}
	return Template{
		htmlTpl: htmlTpl,
	}, nil
}

func ParseFS(fs fs.FS, pattern ...string) (Template, error) {
	htmlTpl := template.New(pattern[0])
	htmlTpl = htmlTpl.Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return "", fmt.Errorf("csrfField not implemented")
			},
			"currentUser": func() (*models.User, error) {
				return nil, fmt.Errorf("currentUser not implemented")
			},
		},
	)
	htmlTpl, err := htmlTpl.ParseFS(fs, pattern...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}
	return Template{
		htmlTpl: htmlTpl,
	}, nil
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}
