package movie

import (
	"log"
	"net/http"

	wsconnection "github.com/RusseLHuang/push-message-example/ws_connection"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type MovieController struct {
	WSConnection *wsconnection.WSConnectionManager
}

func NewMovieController(
	ws *wsconnection.WSConnectionManager,
) MovieController {
	return MovieController{
		WSConnection: ws,
	}
}

func (m MovieController) Send(
	resp http.ResponseWriter,
	req *http.Request,
) {
	vars := mux.Vars(req)
	clientID := vars["clientID"]

	connection := m.WSConnection.GetConnection(clientID)

	// Fire and Forget
	if !m.isConnectionValid(connection) {
		log.Println("Connection is not valid")
		resp.WriteHeader(http.StatusAccepted)
		resp.Write([]byte("Client ID connection is not valid"))
		return
	}

	log.Printf("sending to : %s", vars["clientID"])

	err := connection.WriteMessage(1, []byte("My websocket message"))

	if err != nil {
		log.Println("write:", err)
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("Client ID connection is not valid"))
		return
	}

	resp.WriteHeader(http.StatusOK)
	return
}

func (m MovieController) isConnectionValid(connection *websocket.Conn) bool {
	if connection == nil {
		return false
	}

	log.Println("Trying to read message")
	return true
}
