package message

import (
	"encoding/json"
	"log"
	"strings"
)

// Message structure
type Message struct {
	Header     *Header `json:"headers"`
	Payload    string  `json:"payload"`
	FromBroker bool    `json:"-"`
}

// Create a new message
func Create(payload string) *Message {
	return &Message{
		&Header{},
		payload,
		false,
	}
}

// Text return a new text message
func Text(channel string, payload string) *Message {
	header := &Header{
		"upsub-message-type": TEXT,
		"upsub-channel":      channel,
	}

	return &Message{header, payload, false}
}

// Ping return a new ping message
func Ping() *Message {
	header := &Header{
		"upsub-message-type": PING,
	}

	return &Message{header, "", false}
}

// Pong returns a new pong message
func Pong() *Message {
	header := &Header{
		"upsub-message-type": PONG,
	}

	return &Message{header, "", false}
}

// ResponseAction return a new response action
func ResponseAction(channels []string, action string) *Message {
	for i, channel := range channels {
		channels[i] = channel + ":" + action
	}

	header := &Header{
		"upsub-message-type": TEXT,
		"upsub-channel":      strings.Join(channels, ","),
	}

	return &Message{header, "\"" + strings.Join(channels, ",") + "\"", false}
}

// Decode from byte array to Message
func Decode(message []byte) (*Message, error) {
	var decodedMessage *Message
	err := json.Unmarshal(message, &decodedMessage)
	return decodedMessage, err
}

// Encode Message to a byte array
func Encode(message *Message) ([]byte, error) {
	return json.Marshal(message)
}

// Encode the message instance and return it in an array of bytes
func (m *Message) Encode() ([]byte, error) {
	return Encode(m)
}

// ParseBatch message batch and return an array of message objects
func (m *Message) ParseBatch() []*Message {
	messages := []*Message{}
	err := json.Unmarshal([]byte(m.Payload), &messages)

	if err != nil {
		log.Print("[BATCH INVALID] ", err)
		return messages
	}

	return messages
}
