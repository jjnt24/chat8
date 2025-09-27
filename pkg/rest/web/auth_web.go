package web

import (
	"net/http"

	"github.com/jjnt224/chat8/pkg/view"
)

type AuthWebHandler struct {
	view *view.Renderer
}

func NewAuthWebHandler(v *view.Renderer) *AuthWebHandler {
	return &AuthWebHandler{view: v}
}

func (h *AuthWebHandler) ShowLoginPage(w http.ResponseWriter, r *http.Request) {
	h.view.Render(w, "page-login", nil)
}
