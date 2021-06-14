package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	pushregistry "github.com/RusseLHuang/push-message-example/push_registry"
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
var nodeIP string

func setOutboundIP() {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	nodeIP = localAddr.IP.String()
}

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

	pushregistry.SetPersistentConnectionID(vars["clientID"], nodeIP)

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

	setOutboundIP()
	pushregistry.InitClientConnection()

	clientConnMap = make(map[string]Client)

	router := mux.NewRouter()
	router.HandleFunc("/connect/{clientID}", connect)
	router.HandleFunc("/sendMessage/{clientID}", sendMessage)

	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
