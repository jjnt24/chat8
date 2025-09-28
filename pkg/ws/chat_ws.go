package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jjnt224/chat8/pkg/auth"
	"github.com/jjnt224/chat8/pkg/model"
	"github.com/jmoiron/sqlx"
)

type ChatWebSocket struct {
	DB *sqlx.DB
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type SendMessagePayload struct {
	ReceiverID int64  `json:"receiver_id"`
	Content    string `json:"content"`
}

func (h *ChatWebSocket) ServeWS(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("ws upgrade:", err)
		return
	}
	defer conn.Close()

	for {
		var msg WSMessage
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("ws read:", err)
			break
		}

		switch msg.Type {
		case "send_message":
			var payload SendMessagePayload
			if err := json.Unmarshal(msg.Payload, &payload); err != nil {
				log.Println("invalid payload:", err)
				continue
			}

			// Save to DB
			newMsg := model.Message{
				SenderID:   user.UserID,
				ReceiverID: payload.ReceiverID,
				Room:       " ",
				Content:    payload.Content,
				CreatedAt:  time.Now(),
			}
			_, err := h.DB.NamedExec(`
				INSERT INTO messages (sender_id, receiver_id, room, content, created_at)
				VALUES (:sender_id, :receiver_id, :room, :content, :created_at)
			`, &newMsg)
			if err != nil {
				log.Println("db insert:", err)
				continue
			}

			// Broadcast back to sender
			resp := map[string]interface{}{
				"type": "new_message",
				"payload": map[string]interface{}{
					"id":          newMsg.ID,
					"sender_id":   newMsg.SenderID,
					"receiver_id": newMsg.ReceiverID,
					"content":     newMsg.Content,
					"created_at":  newMsg.CreatedAt,
				},
			}
			conn.WriteJSON(resp)

			// TODO: find receiver connection and send to them
		}
	}
}
