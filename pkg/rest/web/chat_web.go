package web

import (
	"log"
	"net/http"
	"strconv"

	"github.com/jjnt224/chat8/pkg/auth"
	"github.com/jjnt224/chat8/pkg/model"
	"github.com/jjnt224/chat8/pkg/view"

	"github.com/jmoiron/sqlx"
)

type ChatWebHandler struct {
	View *view.Renderer
	DB   *sqlx.DB
}

func (h *ChatWebHandler) ShowChatRoom(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
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

	// Fetch last 20 messages between current user and peer
	var messages []model.Message
	query := `
		SELECT * FROM messages 
		WHERE (sender_id = ? AND receiver_id = ?)
		   OR (sender_id = ? AND receiver_id = ?)
		ORDER BY created_at ASC
		LIMIT 20
	`
	if err := h.DB.Select(&messages, query, user.UserID, peerID, peerID, user.UserID); err != nil {
		log.Panicln(err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	h.View.Render(w, "page-chat", map[string]interface{}{
		"Messages":      messages,
		"CurrentUserID": user.UserID,
		"PeerID":        peerID,
		"PeerUsername":  "Friend", // TODO: fetch real username
	})
}
