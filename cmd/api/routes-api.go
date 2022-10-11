package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	mux.Get("/api/all-users", app.AllUsers)
	mux.Get("/api/getAllMsg/{id}/{uid}", app.GetAllMsg)
	mux.Post("/api/msg/{id}", app.PostMsg)
	mux.Post("/api/add-user", app.AddUser)
	mux.Get("/api/all-friend/{id}", app.GetFriend)
	// mux.Route("/api", func(r chi.Router) {
	// 	mux.Post("/all-users", app.AllUsers)
	// })

	return mux
}
