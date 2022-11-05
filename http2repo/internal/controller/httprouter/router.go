/*
httprouter определяет реакцию на http запросы
*/
package httprouter

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/oauth"

	"github.com/b2b2b-pro/lib/repo_client"
	"go.uber.org/zap"
)

type WebRouter struct {
	R    *chi.Mux
	repo *repo_client.RepoGRPC
	ba   *oauth.BearerAuthentication
}

func (wr *WebRouter) hello(w http.ResponseWriter, r *http.Request) {
	t := r.Header.Get("Authorization")
	ctx := r.Context()
	fmt.Fprintf(w, "Прювет Волку!\n%v\n%v\n%v\n", t, ctx.Value(oauth.CredentialContext), ctx.Value(oauth.AccessTokenContext))

}

func New(db *repo_client.RepoGRPC) *WebRouter {
	zap.S().Debug("Configuring Http router.")
	wr := &WebRouter{}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	wr.ba = oauth.NewBearerAuthentication(
		"b2b2bSecretKey",
		nil)

	r.Use(wr.ba.Authorize)

	r.Get("/", wr.hello)
	r.Mount("/obligation", wr.obligationRouter())
	r.Mount("/entity", wr.entityRouter())
	wr.R = r
	wr.repo = db

	return wr
}
