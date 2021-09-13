package message

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
	"tcp-fragmented-server/internal/pkg/interfaces"
)

type Reader struct {
	reader *bufio.Reader
}

func NewReader(connection net.Conn) *Reader {
	return &Reader{
		reader: bufio.NewReaderSize(connection, 65536),
	}
}

func (r *Reader) Read() (interfaces.Message, error) {
	header, err := r.reader.Peek(2)
	if err != nil {
		return interfaces.Message{}, err
	}

	messageLength := binary.BigEndian.Uint16(header)
	buf := make([]byte, 2+messageLength)

	n, err := io.ReadFull(r.reader, buf)
	if err != nil {
		return interfaces.Message{}, err
	}

	return New(buf[:n]), nil
}
