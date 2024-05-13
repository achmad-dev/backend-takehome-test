package donation

import (
	"github.com/go-chi/chi"
	"net/http"
)

func Routes( /* any dependency injection comes here*/ ) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/{id}", GetDonationByIdHandler)
	r.Post("/create", CreateDonationHandler)
	return r
}

func CreateDonationHandler(rw http.ResponseWriter, r *http.Request) {
	// TODO: implement this
}

func GetDonationByIdHandler(rw http.ResponseWriter, r *http.Request) {
	// TODO: implement this
}
