package message

import (
	"log"
	"strconv"
	"tcp-fragmented-server/internal/pkg/interfaces"
)

const tag = byte('@')

func New(message []byte) interfaces.Message {
	var body []byte

	if len(message) > 2 {
		body = message[2:]
	}

	return interfaces.Message{
		Encoded:   message,
		Body:      body,
		Recipient: decodeTag(body),
	}
}

func decodeTag(body []byte) *interfaces.ClientId {
	var recipient *interfaces.ClientId

	if len(body) < 2 {
		return recipient
	}

	if body[0] != tag {
		return recipient
	}

	var digits []byte

	for _, b := range body[1:] {
		if !isDigit(b) {
			break
		}
		digits = append(digits, b)
	}

	if len(digits) == 0 {
		return recipient
	}

	recipientId, err := strconv.Atoi(string(digits))
	if err != nil {
		log.Printf("Server: can't parse recipient_id '%s'", string(digits))
		return recipient
	}

	return (*interfaces.ClientId)(&recipientId)
}

func isDigit(b byte) bool {
	if b >= '0' && b <= '9' {
		return true
	}

	return false
}
