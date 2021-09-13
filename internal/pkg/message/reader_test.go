package message

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
	"time"
)

func TestReader_Single(t *testing.T) {
	originalMessage := Encode([]byte("this is test message"))
	write, read := net.Pipe()

	reader := NewReader(read)

	go func() {
		_, err := write.Write(originalMessage[:3])
		require.NoError(t, err)

		time.Sleep(100 * time.Millisecond)

		_, err = write.Write(originalMessage[3:])
		require.NoError(t, err)
	}()

	msg, err := reader.Read()
	require.NoError(t, err)
	assert.Equal(t, originalMessage, msg.Encoded)
}

func TestReader_Multiple(t *testing.T) {
	originalMessage1 := Encode([]byte("this is first test message"))
	originalMessage2 := Encode([]byte("this is second test message"))
	write, read := net.Pipe()

	reader := NewReader(read)

	go func() {
		_, err := write.Write(originalMessage1[:3])
		require.NoError(t, err)

		time.Sleep(100 * time.Millisecond)

		_, err = write.Write(originalMessage1[3:])
		require.NoError(t, err)
		_, err = write.Write(originalMessage2[:8])
		require.NoError(t, err)

		time.Sleep(100 * time.Millisecond)

		_, err = write.Write(originalMessage2[8:])
		require.NoError(t, err)
	}()

	msg, err := reader.Read()
	require.NoError(t, err)
	assert.Equal(t, originalMessage1, msg.Encoded)

	msg, err = reader.Read()
	require.NoError(t, err)
	assert.Equal(t, originalMessage2, msg.Encoded)
}
