package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jjnt224/chat8/pkg/auth"
	"github.com/jjnt224/chat8/pkg/model"
	"github.com/jmoiron/sqlx"
)

type ChatAPIHandler struct {
	DB *sqlx.DB
}

// GET /api/messages?peer_id=2&limit=50
func (h *ChatAPIHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	peerIDStr := r.URL.Query().Get("peer_id")
	if peerIDStr == "" {
		http.Error(w, "missing peer_id", http.StatusBadRequest)
		return
	}

	peerID, err := strconv.ParseInt(peerIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid peer_id", http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 200 {
			limit = l
		}
	}

	var messages []model.Message
	query := `
		SELECT * FROM messages 
		WHERE (sender_id = ? AND receiver_id = ?)
		   OR (sender_id = ? AND receiver_id = ?)
		ORDER BY created_at DESC
		LIMIT ?
	`
	if err := h.DB.Select(&messages, query, user.UserID, peerID, peerID, user.UserID, limit); err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	// Reverse messages to ascending order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
