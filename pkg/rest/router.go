package rest

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"

	"github.com/jjnt224/chat8/pkg/auth"
	"github.com/jjnt224/chat8/pkg/config"
	"github.com/jjnt224/chat8/pkg/rest/api"
	"github.com/jjnt224/chat8/pkg/rest/web"
	"github.com/jjnt224/chat8/pkg/view"
)

func NewRouter(cfg config.Config, db *sqlx.DB, store *auth.SessionStore) http.Handler {
	r := chi.NewRouter()
	renderer := view.NewRenderer("pkg/view/templates")

	r.Use(cors.Handler(cors.Options{AllowedOrigins: []string{"*"}, AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, AllowedHeaders: []string{"*"}, MaxAge: 300}))

	// REST endpoints
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })
	// a := NewAuthHandler(db, jwtSvc)
	// m := NewMessageHandler(db)

	// api.NewAuthAPIHandler(db)
	authAPIHandler := &api.AuthAPIHandler{SessionStore: store, DB: db}
	authWebHandler := &web.AuthWebHandler{SessionStore: store, View: renderer}

	// Web routes
	r.Get("/login", authWebHandler.ShowLoginPage)

	r.Get("/register", authWebHandler.ShowRegisterPage)
	r.Group(func(pr chi.Router) { // protected
		pr.Use(auth.AuthMiddleware(store))

		pr.Get("/", authWebHandler.ShowDashboardPage)
	})

	// API routes
	r.Route("/api", func(api chi.Router) {
		api.Post("/register", authAPIHandler.RegisterAPI)
		api.Post("/login", authAPIHandler.LoginAPI)
		api.Post("/logout", authAPIHandler.LogoutAPI)

		api.Group(func(papi chi.Router) { // protected
			papi.Use(auth.AuthMiddleware(store))
			papi.Get("/me", authAPIHandler.MeAPI)
		})
	})

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
