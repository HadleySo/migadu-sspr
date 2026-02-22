package handlers

import (
	"net/http"
	"text/template"

	"golang.hadleyso.com/msspr/src/scenes"
)

func Landing(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFS(scenes.TemplateFS, "scenes/landing.html", "scenes/base.html"))
	tmpl.ExecuteTemplate(w, "base", map[string]any{"PageTitle": "Hello"})
}
