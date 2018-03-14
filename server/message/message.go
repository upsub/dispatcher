package message

import (
	"encoding/json"
	"log"
	"strings"
)

// Message structure
type Message struct {
	Type       string
	Channel    string
	Header     Header
	Payload    string
	FromBroker bool
}

// Create a new message
func Create(payload string) *Message {
	return &Message{
		TEXT,
		"",
		Header{},
		payload,
		false,
	}
}

// Text return a new text message
func Text(channel string, payload string) *Message {
	return &Message{TEXT, channel, Header{}, payload, false}
}

// Json return a new json message
func Json(channel string, payload interface{}) *Message {
	encodedPayload, err := json.Marshal(payload)

	if err != nil {
		log.Print("[ERROR] ", err)
		return nil
	}

	return &Message{JSON, channel, Header{}, string(encodedPayload), false}
}

// Ping return a new ping message
func Ping() *Message {
	return &Message{PING, "", Header{}, "", false}
}

// Pong returns a new pong message
func Pong() *Message {
	return &Message{PONG, "", Header{}, "", false}
}

// ResponseAction return a new response action
func ResponseAction(channels []string, action string) *Message {
	for i, channel := range channels {
		channels[i] = channel + ":" + action
	}

	return &Message{
		TEXT,
		strings.Join(channels, ","),
		Header{},
		strings.Join(channels, ","),
		false,
	}
}

// Decode from byte array to Message
func Decode(message []byte) (*Message, error) {
	return parse(message)
}

// Encode Message to a byte array
func Encode(message *Message) []byte {
	msg := strings.TrimSpace(message.Type + " " + message.Channel)

	for key, value := range message.Header {
		msg += "\n" + key + ": " + value
	}

	if message.Payload != "" {
		msg += "\n\n" + message.Payload
	}

	return []byte(msg)
}

// Encode the message instance and return it in an array of bytes
func (m *Message) Encode() []byte {
	return Encode(m)
}

// ParseBatch message batch and return an array of message objects
func (m *Message) ParseBatch() []*Message {
	messages := []string{}
	err := json.Unmarshal([]byte(m.Payload), &messages)
	decodedMessages := []*Message{}

	if err != nil {
		log.Print("[BATCH INVALID] ", err)
		return decodedMessages
	}

	for _, msg := range messages {
		decoded, _ := Decode([]byte(msg))
		decodedMessages = append(decodedMessages, decoded)
	}

	return decodedMessages
}
