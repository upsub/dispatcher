package dispatcher

import (
	"log"
	"strconv"

	nats "github.com/nats-io/go-nats"
	"github.com/upsub/dispatcher/server/config"
	"github.com/upsub/dispatcher/server/message"
)

type listener func(msg *message.Message)

type broker struct {
	connection *nats.Conn
	config     *config.Config
	listeners  map[string][]listener
}

func createBroker(config *config.Config) *broker {
	if config.Nats == nil {
		return &broker{nil, config, make(map[string][]listener)}
	}

	connection, err := nats.Connect("nats://" + config.Nats.Host + ":" + strconv.Itoa(int(config.Nats.Port)))

	if err != nil {
		log.Print("[ERROR] [NATS] ", err)
		return nil
	}

	broker := &broker{
		connection,
		config,
		make(map[string][]listener),
	}

	broker.connection.Subscribe("upsub.>", func(msg *nats.Msg) {
		decodedMessage, err := message.Decode(msg.Data)
		decodedMessage.FromBroker = true

		if err != nil {
			log.Print("[ERROR] [NATS] ", err)
			return
		}

		broker.emit(msg.Subject, decodedMessage)
	})

	return broker
}

func (b *broker) send(channel string, msg *message.Message) {
	encodedPayload := message.Encode(msg)

	b.connection.Publish(channel, encodedPayload)
}

func (b *broker) emit(channel string, msg *message.Message) *broker {
	for _, callback := range b.listeners[channel] {
		callback(msg)
	}

	return b
}

func (b *broker) on(channel string, callback listener) *broker {
	if b.listeners[channel] != nil {
		b.listeners[channel] = append(b.listeners[channel], callback)
	} else {
		b.listeners[channel] = []listener{callback}
	}

	return b
}
