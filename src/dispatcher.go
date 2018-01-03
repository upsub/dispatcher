package dispatcher

import (
	"log"
	"github.com/upsub/dispatcher/src/util"
	"github.com/upsub/dispatcher/src/message"
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
		case message := <-d.broadcast:
			msg, client := message()
			d.processMessage(msg, client)
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
	msg []byte,
	sender *client,
) {
	dmsg := message.Decode(msg)
	msgType, _ := dmsg.Header["upsub-message-type"]

	switch msgType {
	case message.SubscripeMessage:
		sender.subscribe(dmsg.Header["upsub-channel"])
		break
	case message.UnsubscribeMessage:
		sender.unsubscribe(dmsg.Header["upsub-channel"])
		break
	case message.TextMessage:
	default:
		d.dispatch(
			dmsg,
			map[string]string{"upsub-sender-id": sender.id},
			sender,
		)
	}
}

func (d *dispatcher) dispatch(
	msg *message.Message,
	responseHeaders map[string]string,
	sender *client,
) {
	for _, event :=  range msg.Payload {
		for client := range d.clients {
			if sender != nil && sender == client {
				continue
			}

			if client.appID != sender.appID {
				continue
			}

			if ok := util.Contains(client.subscriptions, event.Channel); !ok {
				continue
			}

			responseMessage, err := message.Encode(
				message.Create(
					util.Merge(responseHeaders, map[string]string{ "type": message.TextMessage }),
					[]*message.Event{event},
				),
			)

			if err != nil {
				log.Print("[FAILED]", err)
				continue
			}

			client.send <- responseMessage
		}
	}
}
