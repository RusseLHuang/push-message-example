package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Client struct {
	clientID string
	conn     *websocket.Conn
}

var clientConnMap map[string]Client

var addr = flag.String("addr", "localhost:8080", "http service address")
var upgrader = websocket.Upgrader{} // use default options

func closeConn(deviceID string) {
	clientConnMap[deviceID].conn.Close()
}

func connect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Print(vars["clientID"])

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	clientConnMap[vars["clientID"]] = Client{
		clientID: vars["clientID"],
		conn:     c,
	}

	defer closeConn(vars["clientID"])

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	conn := clientConnMap[vars["clientID"]].conn

	log.Printf("sending to : %s", vars["clientID"])

	err := conn.WriteMessage(1, []byte("My websocket message"))

	if err != nil {
		log.Println("write:", err)
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	clientConnMap = make(map[string]Client)

	r := mux.NewRouter()
	r.HandleFunc("/connect/{clientID}", connect)
	r.HandleFunc("/sendMessage/{clientID}", sendMessage)

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
