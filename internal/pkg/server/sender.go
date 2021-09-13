package server

import (
	"log"
	"net"
	"tcp-fragmented-server/internal/pkg/interfaces"
)

func (s *Server) sendMessages() {
	for message := range s.newMessages {
		s.sendMessage(message)
	}
}

func (s *Server) sendMessage(message interfaces.Message) {
	var activeConnections []net.Conn

	if message.Recipient == nil { // send to all
		log.Println("Server: no tagged clients, so send to all")
		activeConnections = s.activeClients.getAll()
	} else {
		log.Printf("Server: tagged client #%d\n", *message.Recipient)
		connection, ok := s.activeClients.get(*message.Recipient)
		if !ok {
			log.Printf("Server: client #%d not found (already disconnected?)", *message.Recipient)
		}
		if ok {
			activeConnections = append(activeConnections, connection)
		}
	}

	if len(activeConnections) == 0 {
		log.Println("Server: there are no active clients for this message:", string(message.Body))
		return
	}

	for _, connection := range activeConnections {
		log.Printf("Server: sending message '%s'\n", string(message.Body))
		if _, err := connection.Write(message.Encoded); err != nil {
			log.Printf("Server: failed to send message '%s':%v\n", string(message.Body), err)
			continue
		}
		log.Printf("Server: successfully sent message '%s'\n", string(message.Body))
	}
}
