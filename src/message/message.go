package message

import (
	"encoding/json"
	"log"
)

// Event is an event with in the message payload
type Event struct {
	Name    string
	Channel string
	Body    string
}

// Message is the message structure for communication between server and clients
type Message struct {
	Header  map[string]string
	Payload []*Event
}

// Create a new message with cunstom header and events
func Create(header map[string]string, events []*Event) *Message {
	return &Message{
		Header:  header,
		Payload: events,
	}
}

// Decode message from byte array to Message struct
func Decode(message []byte) *Message {
	var decodedMessage *Message
	err := json.Unmarshal(message, &decodedMessage)

	if err != nil {
		log.Print(err)
	}

	return decodedMessage
}

// Encode Message struct to a byte array and return an error if any was encounterd
func Encode(message *Message) ([]byte, error) {
	return json.Marshal(message)
}
