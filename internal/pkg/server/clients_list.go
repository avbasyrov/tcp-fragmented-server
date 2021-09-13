package server

import (
	"errors"
	"net"
	"sort"
	"sync"
	"tcp-fragmented-server/internal/pkg/interfaces"
)

const maxClients = 256

type clientsList struct {
	clients           []int
	clientConnections map[interfaces.ClientId]net.Conn

	mu sync.RWMutex
}

func newClientsList() *clientsList {
	return &clientsList{
		clientConnections: make(map[interfaces.ClientId]net.Conn),
		mu:                sync.RWMutex{},
	}
}

func (c *clientsList) add(connection net.Conn) (interfaces.ClientId, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.clients) >= maxClients {
		return 0, errors.New("too much clients")
	}

	sort.Ints(c.clients)

	newClientId := 0
	for _, clientId := range c.clients {
		if newClientId != clientId {
			break
		}
		newClientId++
	}

	c.clients = append(c.clients, newClientId)
	c.clientConnections[interfaces.ClientId(newClientId)] = connection

	return interfaces.ClientId(newClientId), nil
}

func (c *clientsList) remove(id interfaces.ClientId) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for idx, clientId := range c.clients {
		if clientId == int(id) {
			c.clients[idx] = c.clients[len(c.clients)-1]
			c.clients = c.clients[:len(c.clients)-1]
			delete(c.clientConnections, id)
			break
		}
	}
}

func (c *clientsList) get(id interfaces.ClientId) (net.Conn, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	connection, ok := c.clientConnections[id]
	return connection, ok
}

func (c *clientsList) getAll() []net.Conn {
	var activeConnections []net.Conn

	c.mu.Lock()
	defer c.mu.Unlock()

	for _, connection := range c.clientConnections {
		activeConnections = append(activeConnections, connection)
	}

	return activeConnections
}
