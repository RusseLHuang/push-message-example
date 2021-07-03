package wsconnection

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type WSConnectionController struct {
	wsConnectionManager *WSConnectionManager
	httpUpgrader        websocket.Upgrader
}

func NewWSConnectionController(
	wsConnectionManager *WSConnectionManager,
) WSConnectionController {
	return WSConnectionController{
		wsConnectionManager: wsConnectionManager,
		httpUpgrader:        websocket.Upgrader{},
	}
}

func (ws WSConnectionController) Connect(
	resp http.ResponseWriter,
	req *http.Request,
) {
	vars := mux.Vars(req)
	clientID := vars["clientID"]
	log.Print(vars["clientID"])

	connection, err := ws.httpUpgrader.Upgrade(resp, req, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	ws.wsConnectionManager.StoreConnection(clientID, connection)

	defer ws.wsConnectionManager.CloseConnection(clientID)

	for {
		mt, message, err := connection.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = connection.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
