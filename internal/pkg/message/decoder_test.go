package message

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"tcp-fragmented-server/internal/pkg/interfaces"
	"testing"
)

func TestTaggedWithComma(t *testing.T) {
	msg := Encode([]byte("@123, hello!"))

	m := New(msg)
	require.NotNil(t, m.Recipient)
	assert.Equal(t, interfaces.ClientId(123), *m.Recipient)
}

func TestTaggedWithSpace(t *testing.T) {
	msg := Encode([]byte("@2 hello!"))

	m := New(msg)
	require.NotNil(t, m.Recipient)
	assert.Equal(t, interfaces.ClientId(2), *m.Recipient)
}

func TestNotTagged(t *testing.T) {
	msg := Encode([]byte("Hello!"))

	m := New(msg)
	assert.Nil(t, m.Recipient)
}
