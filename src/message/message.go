package message

import (
	"encoding/json"
	"log"
)

// Message is the message structure for communication between server and clients
type Message struct {
	Header     *Header `json:"headers"`
	Payload    string  `json:"payload"`
	FromBroker bool    `json:"-"`
}

// Create a new message with cunstom header and events
func Create(payload string) *Message {
	return &Message{
		&Header{},
		payload,
		false,
	}
}

// Text, create a new text message
func Text(channel string, payload string) *Message {
	header := &Header{
		"upsub-message-type": TextMessage,
		"upsub-channel":      channel,
	}

	return &Message{header, payload, false}
}

func (m *Message) Batch() []*Message {
	messages := []*Message{}
	err := json.Unmarshal([]byte(m.Payload), &messages)

	if err != nil {
		log.Print("[BATCH INVALID] ", err)
		return messages
	}

	return messages
}

// Decode message from byte array to Message struct
func Decode(message []byte) (*Message, error) {
	var decodedMessage *Message
	err := json.Unmarshal(message, &decodedMessage)
	return decodedMessage, err
}

// Encode Message struct to a byte array and return an error if any was encounterd
func Encode(message *Message) ([]byte, error) {
	return json.Marshal(message)
}
