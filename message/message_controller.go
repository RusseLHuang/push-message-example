package message

import (
	"log"
	"net/http"

	wsconnection "github.com/RusseLHuang/push-message-example/ws_connection"
	"github.com/gorilla/mux"
)

type MessageController struct {
	WSConnection *wsconnection.WSConnectionManager
}

func NewMessageController(
	ws *wsconnection.WSConnectionManager,
) MessageController {
	return MessageController{
		WSConnection: ws,
	}
}

func (m MessageController) Send(
	resp http.ResponseWriter,
	req *http.Request,
) {
	vars := mux.Vars(req)
	clientID := vars["clientID"]

	connection := m.WSConnection.GetConnection(clientID)

	// Fire and Forget
	if connection == nil {
		resp.WriteHeader(http.StatusAccepted)
		resp.Write([]byte("Client ID connection is not exist in current node"))
		return
	}

	log.Printf("sending to : %s", vars["clientID"])

	err := connection.WriteMessage(1, []byte("My websocket message"))

	if err != nil {
		log.Println("write:", err)
	}

	resp.WriteHeader(http.StatusOK)

}
