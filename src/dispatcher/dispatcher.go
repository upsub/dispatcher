package dispatcher

import (
	"log"

	"github.com/upsub/dispatcher/src/message"
)

type Dispatcher struct {
	connections map[*connection]bool
	broadcast   chan func() ([]byte, *connection)
	register    chan *connection
	unregister  chan *connection
}

func Create() *Dispatcher {
	return &Dispatcher{
		connections: make(map[*connection]bool),
		broadcast:   make(chan func() ([]byte, *connection)),
		register:    make(chan *connection),
		unregister:  make(chan *connection),
	}
}

func (d *Dispatcher) Serve() {
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

			d.processMessage(dmsg, connection)
		}
	}
}

func (d *Dispatcher) connect(connection *connection) {
	log.Println("[Register]", connection.id, connection.name)
	d.connections[connection] = true
}

func (d *Dispatcher) disconnect(connection *connection) {
	log.Println("[Unregister]", connection.id, connection.name)
	if _, ok := d.connections[connection]; ok {
		delete(d.connections, connection)
		close(connection.send)
	}
}

func (d *Dispatcher) processMessage(
	msg *message.Message,
	sender *connection,
) {
	if msg.Header == nil {
		log.Print("[INVALid MESSAGE] ", msg)
		return
	}

	msgType := msg.Header.Get("upsub-message-type")

	switch msgType {
	case message.SubscripeMessage:
		log.Print("[SUBSCRIBE] ", msg.Payload)
		sender.subscribe(msg.Payload)
		break
	case message.UnsubscribeMessage:
		log.Print("[UNSUBSCRIBE] ", msg.Payload)
		sender.unsubscribe(msg.Payload)
		break
	case message.TextMessage:
		responseMessage := message.Create(msg.Payload)
		responseMessage.Header = msg.Header

		d.Dispatch(
			responseMessage,
			sender,
		)
	}
}

func (d *Dispatcher) Dispatch(
	msg *message.Message,
	sender *connection,
) {
	for connection := range d.connections {
		if sender != nil && sender == connection {
			continue
		}

		if sender != nil && connection.appID != sender.appID {
			continue
		}

		if ok := connection.hasSubscription(msg.Header.Get("upsub-channel")); !ok {
			continue
		}

		responseMessage, err := message.Encode(msg)

		if err != nil {
			log.Print("[FAILED]", err)
			continue
		}

		log.Print("[SEND] ", responseMessage)
		connection.send <- responseMessage
	}
}
