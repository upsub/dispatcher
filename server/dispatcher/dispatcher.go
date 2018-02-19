package dispatcher

import (
	"log"
	"strings"

	"github.com/upsub/dispatcher/server/config"
	"github.com/upsub/dispatcher/server/message"
)

type Dispatcher struct {
	broker      *broker
	config      *config.Config
	connections map[*connection]bool
	broadcast   chan func() ([]byte, *connection)
	register    chan *connection
	unregister  chan *connection
}

func Create(config *config.Config) *Dispatcher {
	return &Dispatcher{
		broker:      createBroker(config),
		config:      config,
		connections: make(map[*connection]bool),
		broadcast:   make(chan func() ([]byte, *connection)),
		register:    make(chan *connection),
		unregister:  make(chan *connection),
	}
}

func (d *Dispatcher) Serve() {
	d.broker.on("upsub.dispatcher.message", func(msg *message.Message) {
		d.ProcessMessage(msg, nil)
	})

	for {
		select {
		case connection := <-d.register:
			d.connect(connection)
		case connection := <-d.unregister:
			d.disconnect(connection)
		case payload := <-d.broadcast:
			msg, connection := payload()
			dmsg, err := message.Decode(msg)

			if err != nil {
				log.Print("[MESSAGE DECODE FAILED] ", err)
				continue
			}

			d.ProcessMessage(dmsg, connection)
		}
	}
}

func (d *Dispatcher) connect(connection *connection) {
	log.Println("[REGISTER]", connection.id, connection.name)
	d.connections[connection] = true
}

func (d *Dispatcher) disconnect(connection *connection) {
	log.Println("[UNREGISTER]", connection.id, connection.name)
	if _, ok := d.connections[connection]; ok {
		delete(d.connections, connection)
		close(connection.send)
	}
}

func (d *Dispatcher) ProcessMessage(
	msg *message.Message,
	sender *connection,
) {
	if msg.Header == nil {
		log.Print("[INVALID MESSAGE] ", msg)
		return
	}

	msgType := msg.Header.Get("upsub-message-type")

	switch msgType {
	case message.SUBSCRIBE:
		log.Print("[SUBSCRIBE] ", msg.Payload)
		channels := strings.Split(strings.Replace(msg.Payload, "\"", "", 2), ",")
		sender.subscribe(channels)
		break
	case message.UNSUBSCRIBE:
		log.Print("[UNSUBSCRIBE] ", msg.Payload)
		channels := strings.Split(strings.Replace(msg.Payload, "\"", "", 2), ",")
		sender.unsubscribe(channels)
		break
	case message.PING:
		log.Print("[PING] ", msg.Payload)
		sender.send <- message.Pong()
		break
	case message.BATCH:
		log.Print("[BATCH] ", msg.Payload)
		for _, msg := range msg.ParseBatch() {
			d.ProcessMessage(msg, sender)
		}
		break
	case message.TEXT:
		log.Print("[RECEIVED] ", msg.Header.Get("upsub-channel"), " ", msg.Payload)

		d.Dispatch(
			msg,
			sender,
		)
	}
}

func (d *Dispatcher) Dispatch(
	msg *message.Message,
	sender *connection,
) {
	if !msg.FromBroker {
		d.broker.send("upsub.dispatcher.message", msg)
	}

	for connection := range d.connections {
		if sender != nil && sender == connection {
			continue
		}

		if sender != nil && connection.appID != sender.appID {
			continue
		}

		if !connection.shouldReceive(msg) {
			continue
		}

		if connection.appID != "" {
			msg.Header.Set("upsub-app-id", connection.appID)
		}

		log.Print("[SEND] ", msg.Header.Get("upsub-channel"), " ", msg.Payload)
		connection.send <- msg
	}
}
