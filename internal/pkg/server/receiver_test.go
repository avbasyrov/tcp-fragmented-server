package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"tcp-fragmented-server/internal/pkg/message"
	"testing"
	"time"
)

func TestReceiver(t *testing.T) {
	server := New("")

	serverConnection1, clientConnection1 := net.Pipe()
	go server.receiveMessages(serverConnection1)

	msg := message.Encode([]byte("AB"))

	n, err := clientConnection1.Write(msg)
	require.NoError(t, err)
	require.Equal(t, len(msg), n)

	time.Sleep(100 * time.Millisecond)

	got := <-server.newMessages

	assert.Equal(t, msg, got.Encoded)
}
