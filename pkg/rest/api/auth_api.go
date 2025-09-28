package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/jjnt224/chat8/pkg/auth"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type AuthAPIHandler struct {
	DB           *sqlx.DB
	SessionStore *auth.SessionStore
}

// func NewAuthAPIHandler(db *sqlx.DB) *AuthAPIHandler {
// 	return &AuthAPIHandler{DB: db}
// }

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthAPIHandler) LoginAPI(w http.ResponseWriter, r *http.Request) {
	// Pastikan hanya method POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form values dari request
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	var id, cityID int64
	var hash string
	err := h.DB.QueryRow(`SELECT id, password_hash, city_id FROM users WHERE username = ?`, username).Scan(&id, &hash, &cityID)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		http.Error(w, "invalid creds", http.StatusUnauthorized)
		return
	}

	// Generate random session token
	token, err := auth.GenerateSecureToken(32)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Simpan session ke Redis dengan TTL = 5 menit
	session := auth.SessionData{
		UserID:   1,
		Username: username,
	}
	ctx := context.Background()
	if err := h.SessionStore.Save(ctx, token, session, 5*time.Minute); err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	// Set cookie HttpOnly
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		HttpOnly: true,
		Secure:   true, // wajib diaktifkan jika HTTPS
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  time.Now().Add(5 * time.Minute),
	})

	// Redirect ke halaman utama setelah login sukses
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// type registerRequest struct {
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// 	CityID   int64  `json:"city_id"`
// }

type registerResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	CityID   int64  `json:"city_id"`
}

// RegisterAPI handles user registration
func (h *AuthAPIHandler) RegisterAPI(w http.ResponseWriter, r *http.Request) {
	// Support both JSON and form submissions
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	cityIDStr := r.FormValue("city_id")

	if username == "" || password == "" || cityIDStr == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	cityID, err := strconv.ParseInt(cityIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid City ID", http.StatusBadRequest)
		return
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Insert into DB
	var userID int64
	query := `INSERT INTO users (username, password_hash, city_id) VALUES ($1, $2, $3) RETURNING id`
	err = h.DB.QueryRow(query, username, string(hash), cityID).Scan(&userID)
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(registerResponse{
		ID:       userID,
		Username: username,
		CityID:   cityID,
	})
}

// LogoutAPI menghapus session dari Redis dan cookie
func (h *AuthAPIHandler) LogoutAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Ambil session token dari cookie
	cookie, err := r.Cookie("session_token")
	if err == nil {
		// Hapus session di Redis
		ctx := context.Background()
		_ = h.SessionStore.Delete(ctx, cookie.Value)

		// Clear cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		})
	}

	// Redirect ke halaman login
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *AuthAPIHandler) MeAPI(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"user":"` + user.Username + `"}`))
}
