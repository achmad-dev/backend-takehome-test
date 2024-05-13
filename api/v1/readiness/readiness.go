package readiness

import (
	"github.com/go-chi/chi"
	"github.com/kitabisa/backend-takehome-test/internal/config"
	"net/http"
)

func Routes( /* any dependency injection comes here*/ ) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/db", ReadinessHandler)
	return r
}

func ReadinessHandler(rw http.ResponseWriter, r *http.Request) {
	var dbConn config.DbConnection
	err := dbConn.GetDbConnectionPool().Db.Ping()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(200)
}
