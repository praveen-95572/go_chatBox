package main

import (
	"fmt"
	"net/http"

	"chat/internal/models"

	"github.com/gorilla/websocket"
)

type WebSocketConnection struct {
	*websocket.Conn
}

type WsPayload struct {
	Action   string              `json:"action"`
	UserID   int                 `json:"user_id"`
	Conn     WebSocketConnection `json:"-"`
	Message  string              `json:"msg"`
	InsertID int
}

type WsJsonResponse struct {
	Chat    []*models.Chat `json:"chat"`
	Action  string         `json:"action"`
	Message string         `json:"message"`
	UserID  int            `json:"user_id"`
}

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var clients = make(map[WebSocketConnection]string)

var wsChan = make(chan WsPayload)

func (app *application) WsEndPoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	app.infoLog.Println(fmt.Sprintf("Client connected from %s", r.RemoteAddr))
	var response WsJsonResponse
	response.Message = "Connected to inoknin"

	err = ws.WriteJSON(response)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	conn := WebSocketConnection{Conn: ws}
	clients[conn] = ""

	go app.ListenForWs(&conn)
}

func (app *application) ListenForWs(conn *WebSocketConnection) { // receive msg
	defer func() {
		if r := recover(); r != nil {
			app.errorLog.Println("ERORR:", fmt.Sprintf("%v", r))
		}
	}()
	var payload WsPayload
	for {
		err := conn.ReadJSON(&payload)
		if err != nil {

		} else if payload.Action == "POST" {
			payload.Conn = *conn
			id, _ := app.DB.PostMsg(payload.UserID, payload.Message)
			payload.InsertID = id
			wsChan <- payload
		} else if payload.Action == "GET" {
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}

func (app *application) ListenToWsChannel() {
	var response WsJsonResponse
	for {
		e := <-wsChan
		response.UserID = e.UserID
		if e.Action == "GET" {
			response.Chat, _ = app.DB.GetAllMsg(e.UserID)
			response.Action = "GET"
		} else if e.Action == "POST" {
			response.Chat, _ = app.DB.GetMsg(e.InsertID)
			response.Action = "POST"
		}

		app.broadcastToAll(response)
	}
}

func (app *application) broadcastToAll(response WsJsonResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			app.errorLog.Printf("Websocket err on %s: %s", response.Action, err)
			_ = client.Close()
			delete(clients, client)
		}
	}
}
