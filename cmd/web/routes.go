package main

import (
	"net/http"

	chi "github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Get("/", app.Home)
	mux.Get("/ws", app.WsEndPoint)
	mux.Get("/person/A", app.personA)
	mux.Get("/user/{id}", app.OneUser)

	return mux
}
