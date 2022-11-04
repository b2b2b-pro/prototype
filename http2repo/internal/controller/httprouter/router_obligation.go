package httprouter

import (
	"fmt"
	"net/http"

	"github.com/b2b2b-pro/lib/object"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/oauth"
	"go.uber.org/zap"
)

func (wr *WebRouter) obligationRouter() chi.Router {
	zap.S().Debug("Configuring Obligation Http router.")

	r := chi.NewRouter()
	r.Get("/", wr.listObligation)
	r.Post("/", wr.createObligation)
	r.Route("/{obligationID}", func(r chi.Router) {
		r.Get("/", getObligation)
	})

	return r
}

// curl http://localhost:8088/obligation
func (wr *WebRouter) listObligation(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Список требований:\n")

	ctx := r.Context()
	tmp, err := wr.repo.ListObligation(ctx.Value(oauth.AccessTokenContext).(string))
	if err != nil {
		zap.S().Error("Repo ListObligation error: ", err, "\n")
	}

	fmt.Fprintf(w, "%v\n", tmp)
}

// curl -X POST -d '{"debtor_id":3,"creditor_id":4, "cost": 20, "origin_id":1, "payment_date":"2020-10-10"}' http://localhost:8088/obligation
// TODO не доделана передача в gRPC
func (wr *WebRouter) createObligation(w http.ResponseWriter, r *http.Request) {
	var err error

	fmt.Fprintf(w, "Поступило требование:\n%v\n", r.Body)

	l, err := object.ParseObligation(r.Body)
	if err != nil {
		zap.S().Error("Error ParseObligation: ", err, "\n")
		return
	}

	zap.S().Debug("CreateObligation получил ", l, " от ParseObligation\n")

	ctx := r.Context()
	l.ID, err = wr.repo.CreateObligation(ctx.Value(oauth.AccessTokenContext).(string), *l)
	if err != nil {
		zap.S().Error("wr.db.Create error: ", err, "\n")
	}
}

// TODO curl http://localhost:8088/obligation/55
func getObligation(w http.ResponseWriter, r *http.Request) {
	obligationID := chi.URLParam(r, "obligationID")
	fmt.Fprintf(w, "GetObligation: %v\n", obligationID)
}
