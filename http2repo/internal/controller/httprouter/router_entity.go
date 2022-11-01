package httprouter

import (
	"fmt"
	"net/http"

	"github.com/b2b2b-pro/lib/object"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/oauth"
	"go.uber.org/zap"
)

func (wr *WebRouter) entityRouter() chi.Router {
	zap.S().Debug("Configuring Entity Http router.")

	r := chi.NewRouter()
	r.Get("/", wr.listEntity)
	r.Post("/", wr.createEntity)
	r.Route("/{entityINN}", func(r chi.Router) {
		r.Get("/", getEntity)
	})

	return r
}

// curl http://localhost:8088/entity
func (wr *WebRouter) listEntity(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Список фирм:\n")

	ctx := r.Context()
	tmp, err := wr.repo.ListEntity(ctx.Value(oauth.AccessTokenContext).(string))
	/*
		cannot use ctx.Value(oauth.AccessTokenContext) (value of type any) as string value in argument to wr.repo.ListEntity: need type assertion
	*/
	if err != nil {
		zap.S().Error("Repo ListEntity error: ", err, "\n")
	}

	fmt.Fprintf(w, "%v\n", tmp)
}

// curl -X POST -d '{"entity_inn":"1000", "entity_kpp":"1000", "short_name":"тест"}' http://localhost:8088/entity
func (wr *WebRouter) createEntity(w http.ResponseWriter, r *http.Request) {
	var err error

	fmt.Fprintf(w, "Поступила информация по фирме:\n%v\n", r.Body)

	frm, err := object.ParseEntity(r.Body)
	if err != nil {
		zap.S().Error("Error newEntity: ", err, "\n")
		return
	}

	zap.S().Info("createEntity получил ", frm, " от ParseEntity\n")

	ctx := r.Context()
	frm.ID, err = wr.repo.CreateEntity(ctx.Value(oauth.AccessTokenContext).(string), *frm)
	if err != nil {
		zap.S().Error("wr.db.Create error: ", err, "\n")
	}
}

// TODO curl http://localhost:8088/entity/55
func getEntity(w http.ResponseWriter, r *http.Request) {
	entityINN := chi.URLParam(r, "entityINN")
	fmt.Fprintf(w, "getEntity: %v\n", entityINN)
}
