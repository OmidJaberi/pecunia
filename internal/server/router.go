package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Routes(api *API) http.Handler {
	r := chi.NewRouter()

	r.Post("/users", api.CreateUser)
	r.Get("/users/{id}/assets", api.ListAssets)
	r.Post("/users/{id}/assets", api.CreateAsset)

	// health
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	return r
}
