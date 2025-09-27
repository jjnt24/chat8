package rest

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"

	"github.com/jjnt224/chat8/pkg/config"
	"github.com/jjnt224/chat8/pkg/rest/web"
	"github.com/jjnt224/chat8/pkg/view"
)

func NewRouter(cfg config.Config, db *sqlx.DB) http.Handler {
	r := chi.NewRouter()
	renderer := view.NewRenderer("pkg/view/templates")

	r.Use(cors.Handler(cors.Options{AllowedOrigins: []string{"*"}, AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, AllowedHeaders: []string{"*"}, MaxAge: 300}))

	// jwtSvc := auth.JWTService{
	// 	AccessSecret:  []byte(cfg.JWTAccessSecret),
	// 	RefreshSecret: []byte(cfg.JWTRefreshSecret),
	// 	AccessTTL:     time.Duration(cfg.JWTAccessTTLMin) * time.Minute,
	// 	RefreshTTL:    time.Duration(cfg.JWTRefreshTTLDays) * 24 * time.Hour,
	// }

	// REST endpoints
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })
	// a := NewAuthHandler(db, jwtSvc)
	// m := NewMessageHandler(db)

	authWebHandler := web.NewAuthWebHandler(renderer)

	r.Get("/login", authWebHandler.ShowLoginPage)

	// r.Post("/register", a.Register)
	// r.Post("/api/login", a.Login)

	// r.Post("/login", a.LoginForm)

	// r.Post("/refresh", a.Refresh)

	// r.Get("/login", a.ShowLoginPage)
	// r.Post("/logout", a.Logout)

	// r.Group(func(pr chi.Router) {
	// 	pr.Use(a.CookieAuth)

	// 	pr.Get("/", a.ShowDashboardPage)

	// })

	// // WebSocket hub and endpoint
	// hub := ws.NewHub()
	// go hub.Run()
	// r.Get("/ws", ws.MakeWSHandler(hub, a.ParseUserFromRequest))

	// === Serve folder static ===
	workDir, _ := os.Getwd() // direktori kerja sekarang
	filesDir := filepath.Join(workDir, "static")

	// Semua file di /static bisa diakses langsung
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir(filesDir))))

	return r
}
