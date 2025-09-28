package model

import "time"

type Message struct {
	ID         int64     `db:"id" json:"id"`
	SenderID   int64     `db:"sender_id" json:"sender_id"`
	ReceiverID int64     `db:"receiver_id" json:"receiver_id"`
	Room       string    `db:"room" json:"room"`
	Content    string    `db:"content" json:"content"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}
