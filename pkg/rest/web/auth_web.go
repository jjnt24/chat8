package web

import (
	"net/http"

	"github.com/jjnt224/chat8/pkg/auth"
	"github.com/jjnt224/chat8/pkg/view"
)

type AuthWebHandler struct {
	View         *view.Renderer
	SessionStore *auth.SessionStore
}

// func NewAuthWebHandler(v *view.Renderer) *AuthWebHandler {
// 	return &AuthWebHandler{View: v}
// }

func (h *AuthWebHandler) ShowLoginPage(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err == nil && cookie.Value != "" {
		// Check if token exists in Redis
		session, _ := h.SessionStore.Get(r.Context(), cookie.Value)
		if session != nil {
			// Token valid → redirect to dashboard
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	// Token invalid or missing → show login page
	h.View.Render(w, "page-login", nil)
}

func (h *AuthWebHandler) ShowDashboardPage(w http.ResponseWriter, r *http.Request) {
	h.View.Render(w, "page-dashboard", nil)
}

func (h *AuthWebHandler) ShowRegisterPage(w http.ResponseWriter, r *http.Request) {
	h.View.Render(w, "page-register", nil)
}
