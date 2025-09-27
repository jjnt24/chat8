package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type AuthAPIHandler struct {
	DB *sqlx.DB
}

func NewAuthAPIHandler(db *sqlx.DB) *AuthAPIHandler {
	return &AuthAPIHandler{DB: db}
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
