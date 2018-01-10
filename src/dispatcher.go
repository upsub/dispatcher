package dispatcher

import (
	"log"

	"github.com/upsub/dispatcher/src/message"
	"github.com/upsub/dispatcher/src/util"
)

type dispatcher struct {
	clients    map[*client]bool
	broadcast  chan func() ([]byte, *client)
	register   chan *client
	unregister chan *client
}

func createDispatcher() *dispatcher {
	return &dispatcher{
		clients:    make(map[*client]bool),
		broadcast:  make(chan func() ([]byte, *client)),
		register:   make(chan *client),
		unregister: make(chan *client),
	}
}

func (d *dispatcher) serve() {
	for {
		select {
		case client := <-d.register:
			d.connect(client)
		case client := <-d.unregister:
			d.disconnect(client)
		case payload := <-d.broadcast:
			msg, client := payload()
			dmsg, err := message.Decode(msg)

			if err != nil {
				log.Print("[MESSAGE DECODE FAILED] ", err)
				continue
			}

			d.processMessage(dmsg, client)
		}
	}
}

func (d *dispatcher) connect(client *client) {
	log.Println("[Register]", client.id, client.name)
	d.clients[client] = true
}

func (d *dispatcher) disconnect(client *client) {
	log.Println("[Unregister]", client.id, client.name)
	if _, ok := d.clients[client]; ok {
		delete(d.clients, client)
		close(client.send)
	}
}

func (d *dispatcher) processMessage(
	msg *message.Message,
	sender *client,
) {
	if msg.Header == nil {
		log.Print("[INVALID MESSAGE] ", msg)
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
	default:
		responseMessage := message.Create(msg.Payload)
		responseMessage.Header = msg.Header

		d.dispatch(
			responseMessage,
			sender,
		)
	}
}

func (d *dispatcher) dispatch(
	msg *message.Message,
	sender *client,
) {
	for client := range d.clients {
		if sender != nil && sender == client {
			continue
		}

		if sender != nil && client.appID != sender.appID {
			continue
		}

		if ok := util.Contains(client.subscriptions, msg.Header.Get("upsub-channel")); !ok {
			continue
		}

		responseMessage, err := message.Encode(msg)

		if err != nil {
			log.Print("[FAILED]", err)
			continue
		}

		log.Print("[SEND] ", responseMessage)
		client.send <- responseMessage
	}
}
