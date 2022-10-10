package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *application) AllUsers(w http.ResponseWriter, r *http.Request) {

	allUsers, err := app.DB.GetAllUser()

	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	app.infoLog.Println("------------------ ALL users extracted ------------------------")

	app.writeJSON(w, http.StatusOK, allUsers)
}

func (app *application) GetAllMsg(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, _ := strconv.Atoi(id)

	chats, err := app.DB.GetAllMsg(userID)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	app.writeJSON(w, http.StatusOK, chats)
}

func (app *application) PostMsg(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, _ := strconv.Atoi(id)

	var userInput struct {
		Message string `json:"msg"`
	}

	err := app.readJSON(w, r, &userInput)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	_, err = app.DB.PostMsg(userID, userInput.Message)

	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	_ = app.writeJSON(w, http.StatusOK, nil)
}

func (app *application) AddUser(w http.ResponseWriter, r *http.Request) {
	var userInput struct {
		Name     string `json:"name"`
		Username string `json:"u_name"`
	}

	err := app.readJSON(w, r, &userInput)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	fmt.Println("NEW USER ", userInput)
	_, err = app.DB.AddNewUser(userInput.Name, userInput.Username)

	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	_ = app.writeJSON(w, http.StatusOK, nil)
}
