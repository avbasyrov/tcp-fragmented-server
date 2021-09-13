package interfaces

type Message struct {
	Encoded   []byte
	Body      []byte
	Recipient *ClientId
}
