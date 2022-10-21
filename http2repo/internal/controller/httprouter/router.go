/*
httprouter определяет реакцию на http запросы
*/
package httprouter

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/b2b2b-pro/lib/repo_client"
	"go.uber.org/zap"
)

type WebRouter struct {
	R    *chi.Mux
	repo *repo_client.RepoGRPC
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Прювет Волку!\n")
}

func New(db *repo_client.RepoGRPC) *WebRouter {
	zap.S().Debug("Configuring Http router.")
	wr := &WebRouter{}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", hello)
	r.Mount("/obligation", wr.obligationRouter())
	r.Mount("/entity", wr.entityRouter())
	wr.R = r
	wr.repo = db

	return wr
}
