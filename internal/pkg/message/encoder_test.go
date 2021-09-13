package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncode(t *testing.T) {
	msg := Encode([]byte("AB"))
	assert.Equal(t, []byte{0x00, 0x02, 0x41, 0x42}, msg)
}
