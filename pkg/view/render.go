package view

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type Renderer struct {
	templates map[string]*template.Template
}

func NewRenderer(templateDir string) *Renderer {
	templates := make(map[string]*template.Template)

	for _, value := range []string{"login", "register", "dashboard", "chat"} {
		templates["page-"+value] = template.Must(template.ParseFiles(
			filepath.Join(templateDir, "base.html"),
			filepath.Join(templateDir, value+".html"),
		))
	}

	return &Renderer{templates: templates}
}

func (r *Renderer) Render(w http.ResponseWriter, name string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, ok := r.templates[name]
	if !ok {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
	err := tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
