package view

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type Renderer struct {
	templates *template.Template
}

func NewRenderer(templateDir string) *Renderer {
	// Parse semua file HTML di folder templates
	templates, err := template.ParseGlob(filepath.Join(templateDir, "*.html"))
	if err != nil {
		panic("failed to parse templates: " + err.Error())
	}
	return &Renderer{templates: templates}
}

func (r *Renderer) Render(w http.ResponseWriter, name string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := r.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
