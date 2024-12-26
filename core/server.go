package core

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type Server struct {
	addr        string
	listener    net.Listener
	connections map[string]net.Conn
	quitChannel chan struct{}
}

func NewServer(addr string) *Server {
	return &Server{
		addr:        addr,
		connections: make(map[string]net.Conn),
		quitChannel: make(chan struct{}),
	}
}

func (server *Server) acceptConnection() {
	for {
		connection, err := server.listener.Accept()
		if err != nil {
			log.Fatal("Failed to accept connection: ", err)
			break
		}

		go server.handleConnection(connection)
	}
}

func (server *Server) handleConnection(connection net.Conn) {
	defer server.closeConnection(connection)

	connectionAddr := connection.RemoteAddr().String()
	server.connections[connectionAddr] = connection

	fmt.Printf("Client connected: %s\n", connectionAddr)

	buffer := make([]byte, 4096)
	for {
		n, err := connection.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Client disconnected: %s\n", connection.RemoteAddr().String())
			} else {
				fmt.Println("Failed to receiving message: ", err)
			}

			break
		}

		server.broadcastMessage(connection, buffer[:n])
	}
}

func (server *Server) broadcastMessage(currentConnection net.Conn, message []byte) {
	currentConnectionAddr := currentConnection.RemoteAddr().String()
	formattedMessage := fmt.Sprintf("[%s] %s", currentConnectionAddr, strings.TrimSpace(string(message)))

	for addr, connection := range server.connections {
		if addr != currentConnectionAddr {
			connection.Write([]byte(formattedMessage))
		}
	}
}

func (server *Server) closeConnection(connection net.Conn) {
	defer connection.Close()
	delete(server.connections, connection.RemoteAddr().String())
}

func (server *Server) Run(callback ...func()) error {
	listener, err := net.Listen("tcp", server.addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	server.listener = listener
	go server.acceptConnection()

	if len(callback) > 0 {
		callback[0]()
	}

	<-server.quitChannel

	return nil
}
