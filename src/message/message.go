package message

import (
	"log"
	"encoding/json"
)

type Event struct {
	Name    string
	Channel string
	Body    string
}

type Message struct {
	Header  map[string]string
	Payload []*Event
}

func Create(header map[string]string, events []*Event) *Message {
	return &Message{
		Header: header,
		Payload: events,
	}
}

func Decode(message []byte) *Message {
	var decodedMessage *Message
	err := json.Unmarshal(message, &decodedMessage)

	if err != nil {
		log.Print(err)
	}

	return decodedMessage
}

func Encode(message *Message) ([]byte, error) {
	return json.Marshal(message)
}
