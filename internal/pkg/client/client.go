package client

import (
	"io"
	"log"
	"net"
	"tcp-fragmented-server/internal/pkg/message"
)

type Client struct {
	connection net.Conn
	address    string
	name       string
}

func New(address string, name string) *Client {
	return &Client{
		address: address,
		name:    name,
	}
}

func (c *Client) Connect() error {
	connection, err := net.Dial("tcp", c.address)
	if err != nil {
		return err
	}

	c.connection = connection

	return nil
}

func (c *Client) SendMessage(msg string) error {
	_, err := c.connection.Write(message.Encode([]byte(msg)))
	return err
}

func (c *Client) Run() error {
	reader := message.NewReader(c.connection)

	for {
		msg, err := reader.Read()

		switch err {
		case nil:
			log.Printf("Client (%s): message received from server: %s\n", c.name, string(msg.Body))
		case io.EOF:
			log.Printf("Client (%s): server closed the connection\n", c.name)
			return nil
		default:
			log.Printf("Client (%s): reader error: %v\n", c.name, err)
			return err
		}
	}
}
