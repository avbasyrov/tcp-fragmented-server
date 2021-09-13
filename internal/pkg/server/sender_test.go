package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"sync"
	"tcp-fragmented-server/internal/pkg/message"
	"testing"
	"time"
)

func TestSendMessage_ToAll(t *testing.T) {
	server := New("")

	serverConnection1, clientConnection1 := net.Pipe()
	go server.receiveMessages(serverConnection1)
	serverConnection2, clientConnection2 := net.Pipe()
	go server.receiveMessages(serverConnection2)

	go server.sendMessages()

	msg := message.Encode([]byte("AB"))
	server.newMessages <- message.New(msg)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		readByClient1 := make([]byte, 100)
		_ = clientConnection1.SetReadDeadline(time.Now().Add(time.Second))
		n, err := clientConnection1.Read(readByClient1)
		require.NoError(t, err)
		assert.Equal(t, 4, n)
		assert.Equal(t, msg, readByClient1[:n])
	}()

	go func() {
		defer wg.Done()
		readByClient2 := make([]byte, 100)
		_ = clientConnection2.SetReadDeadline(time.Now().Add(time.Second))
		n, err := clientConnection2.Read(readByClient2)
		require.NoError(t, err)
		assert.Equal(t, 4, n)
		assert.Equal(t, msg, readByClient2[:n])
	}()

	wg.Wait()
}

func TestSendMessage_Tagged(t *testing.T) {
	server := New("")

	serverConnection1, clientConnection1 := net.Pipe()
	go server.receiveMessages(serverConnection1)

	time.Sleep(100 * time.Millisecond)

	serverConnection2, clientConnection2 := net.Pipe()
	go server.receiveMessages(serverConnection2)

	go server.sendMessages()

	msg := message.Encode([]byte("@1, hello!"))
	server.newMessages <- message.New(msg)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		readByClient1 := make([]byte, 100)
		_ = clientConnection1.SetReadDeadline(time.Now().Add(time.Second))
		n, err := clientConnection1.Read(readByClient1)
		require.Error(t, err)
		assert.True(t, err.(net.Error).Timeout())
		assert.Equal(t, 0, n)
	}()

	go func() {
		defer wg.Done()
		readByClient2 := make([]byte, 100)
		_ = clientConnection2.SetReadDeadline(time.Now().Add(time.Second))
		n, err := clientConnection2.Read(readByClient2)
		require.NoError(t, err)
		assert.Equal(t, 12, n)
		assert.Equal(t, msg, readByClient2[:n])
	}()

	wg.Wait()
}

func TestSendMessage_TaggedNotExisting(t *testing.T) {
	server := New("")

	serverConnection1, clientConnection1 := net.Pipe()
	go server.receiveMessages(serverConnection1)
	serverConnection2, clientConnection2 := net.Pipe()
	go server.receiveMessages(serverConnection2)

	go server.sendMessages()

	msg := message.Encode([]byte("@2, hello!"))
	server.newMessages <- message.New(msg)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		readByClient1 := make([]byte, 100)
		_ = clientConnection1.SetReadDeadline(time.Now().Add(time.Second))
		n, err := clientConnection1.Read(readByClient1)
		require.Error(t, err)
		assert.True(t, err.(net.Error).Timeout())
		assert.Equal(t, 0, n)
	}()

	go func() {
		defer wg.Done()
		readByClient2 := make([]byte, 100)
		_ = clientConnection2.SetReadDeadline(time.Now().Add(time.Second))
		n, err := clientConnection2.Read(readByClient2)
		require.Error(t, err)
		assert.True(t, err.(net.Error).Timeout())
		assert.Equal(t, 0, n)
	}()

	wg.Wait()
}
