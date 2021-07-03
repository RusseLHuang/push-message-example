package node

import (
	"log"
	"net"
)

type Node struct {
	endpoint string
}

func NewNode() *Node {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	nodeIP := localAddr.IP.String()
	log.Println("Node IP: ", nodeIP)

	return &Node{
		endpoint: nodeIP,
	}
}

func (n *Node) GetEndpoint() string {
	return n.endpoint
}
