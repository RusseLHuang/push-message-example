package wsconnection

import (
	"github.com/RusseLHuang/push-message-example/node"
	pushregistry "github.com/RusseLHuang/push-message-example/push_registry"
	"github.com/gorilla/websocket"
)

type WSConnectionManager struct {
	connection  map[string]WSClient
	registry    *pushregistry.PushRegistry
	nodeService *node.Node
}

type WSClient struct {
	clientID string
	conn     *websocket.Conn
}

func NewWSConnectionManager(
	registry *pushregistry.PushRegistry,
	nodeService *node.Node,
) *WSConnectionManager {
	return &WSConnectionManager{
		connection:  make(map[string]WSClient),
		registry:    registry,
		nodeService: nodeService,
	}
}

func (ws *WSConnectionManager) StoreConnection(key string, connection *websocket.Conn) {
	ws.connection[key] = WSClient{
		clientID: key,
		conn:     connection,
	}

	ws.registry.SetPersistentConnectionID(key, ws.nodeService.GetEndpoint())
}

func (ws *WSConnectionManager) GetConnection(key string) *websocket.Conn {
	return ws.connection[key].conn
}

func (ws *WSConnectionManager) CloseConnection(key string) {
	ws.connection[key].conn.Close()
}
