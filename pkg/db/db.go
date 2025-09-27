package db

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	// _ "github.com/jackc/pgx/v5/stdlib" // enable if using pgx as driver

	"github.com/jjnt224/chat8/pkg/config"
)

type DB = sqlx.DB

func MustInit(cfg config.Config) *sqlx.DB {
	db, err := sqlx.Open(cfg.DBDriver, cfg.DBDSN)
	if err != nil {
		log.Fatalln("db open:", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalln("db ping:", err)
	}

	// Pragmas for SQLite only
	if strings.Contains(strings.ToLower(cfg.DBDriver), "sqlite") {
		if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
			log.Println("warn pragma:", err)
		}
	}

	migrate(db)
	return db
}

func migrate(db *sqlx.DB) {
	// users
	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS users (
id INTEGER PRIMARY KEY AUTOINCREMENT,
username TEXT UNIQUE NOT NULL,
password_hash TEXT NOT NULL,
city_id INTEGER NOT NULL,
created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
)`)
	// messages
	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS messages (
id INTEGER PRIMARY KEY AUTOINCREMENT,
sender_id INTEGER NOT NULL,
receiver_id INTEGER,
room TEXT,
content TEXT NOT NULL,
created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
FOREIGN KEY(sender_id) REFERENCES users(id)
)`)
	// refresh tokens (server-tracked, optional)
	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS refresh_tokens (
id INTEGER PRIMARY KEY AUTOINCREMENT,
user_id INTEGER NOT NULL,
token TEXT NOT NULL,
expires_at TIMESTAMP NOT NULL,
created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
FOREIGN KEY(user_id) REFERENCES users(id)
)`)
}
