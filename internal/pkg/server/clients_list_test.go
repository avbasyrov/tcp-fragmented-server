package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClientsList(t *testing.T) {
	c := newClientsList()

	for i := 0; i < maxClients; i++ {
		clientId, err := c.add(nil)
		require.NoError(t, err)
		assert.Equal(t, i, int(clientId))

		_, ok := c.get(clientId)
		assert.True(t, ok)
	}

	assert.Equal(t, maxClients, len(c.getAll()))

	_, err := c.add(nil)
	assert.Error(t, err)

	c.remove(5)
	_, ok := c.get(5)
	assert.False(t, ok)

	clientId, err := c.add(nil)
	require.NoError(t, err)
	assert.Equal(t, 5, int(clientId))

	_, ok = c.get(5)
	assert.True(t, ok)
}
