package server

import (
	"log"
	"net"
	"tcp-fragmented-server/internal/pkg/interfaces"
)

type Server struct {
	activeClients *clientsList
	address       string
	newMessages   chan interfaces.Message
}

func New(address string) *Server {
	return &Server{
		activeClients: newClientsList(),
		address:       address,
		newMessages:   make(chan interfaces.Message, 50),
	}
}

func (s *Server) Serve() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	defer func() {
		err := listener.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	go s.sendMessages()

	for {
		con, err := listener.Accept()
		if err != nil {
			log.Printf("Server: can't accept connection: %+v", err)
			continue
		}

		go s.receiveMessages(con)
	}
}
