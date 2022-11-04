/*
httprouter определяет реакцию на http запросы
*/
package httprouter

import (
	"time"

	"github.com/b2b2b-pro/prototype/auth/internal/controller/client"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/oauth"

	"go.uber.org/zap"
)

type WebRouter struct {
	R *chi.Mux
}

func New() *WebRouter {
	zap.S().Debug("Configuring Auth Http router.")
	wr := &WebRouter{}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE", "HEAD", "OPTION"},
		AllowedHeaders:   []string{"User-Agent", "Content-Type", "Accept", "Accept-Encoding", "Accept-Language", "Cache-Control", "Connection", "DNT", "Host", "Origin", "Pragma", "Referer"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	registerAPI(r)
	wr.R = r
	return wr
}

func registerAPI(r *chi.Mux) {
	s := oauth.NewBearerServer(
		"b2b2bSecretKey",
		time.Second*120,
		&client.Verifier{},
		nil)
	r.Post("/token", s.UserCredentials)
	r.Post("/auth", s.ClientCredentials)
}
