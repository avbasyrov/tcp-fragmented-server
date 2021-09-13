package message

import "encoding/binary"

func Encode(text []byte) []byte {
	length := uint16(len(text))

	msg := make([]byte, 2)
	binary.BigEndian.PutUint16(msg, length)
	msg = append(msg, text...)

	return msg
}
