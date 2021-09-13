package server

import (
	"io"
	"log"
	"net"
	"tcp-fragmented-server/internal/pkg/message"
)

func (s *Server) receiveMessages(connection net.Conn) {
	clientId, err := s.activeClients.add(connection)
	if err != nil {
		log.Println("Server: can't add client:", err)
		return
	}
	defer func() {
		err := connection.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	log.Printf("Server: new client connected #%d\n", clientId)

	reader := message.NewReader(connection)

	for {
		msg, err := reader.Read()

		switch err {
		case nil:
			log.Printf("Server: message received from client: '%s'\n", string(msg.Body))
			s.newMessages <- msg
		case io.EOF:
			log.Println("Server: client closed the connection")
			return
		default:
			log.Printf("Server: reader error: %v\n", err)
			return
		}
	}
}
