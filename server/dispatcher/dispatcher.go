package dispatcher

import (
	"log"
	"strings"

	"github.com/upsub/dispatcher/server/auth"
	"github.com/upsub/dispatcher/server/config"
	"github.com/upsub/dispatcher/server/message"
	"github.com/upsub/dispatcher/server/util"
)

var reservedChannels = []string{
	"upsub/auth/create",
	"upsub/auth/update",
	"upsub/auth/delete",
}

// Dispatcher type
type Dispatcher struct {
	broker      *broker
	config      *config.Config
	store       *auth.Store
	connections map[*connection]bool
	broadcast   chan func() (*message.Message, *connection)
	register    chan *connection
	unregister  chan *connection
}

// Create returns a new instance of the Dispatcher
func Create(config *config.Config, store *auth.Store) *Dispatcher {
	return &Dispatcher{
		broker:      createBroker(config),
		config:      config,
		store:       store,
		connections: make(map[*connection]bool),
		broadcast:   make(chan func() (*message.Message, *connection)),
		register:    make(chan *connection),
		unregister:  make(chan *connection),
	}
}

// Serve listens for incomming connections and messages
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
			d.ProcessMessage(msg, connection)
		}
	}
}

func (d *Dispatcher) connect(connection *connection) {
	log.Println("[REGISTER]", connection.id, connection.name)
	d.connections[connection] = true
}

func (d *Dispatcher) disconnect(conn *connection) {
	log.Println("[UNREGISTER]", conn.id, conn.name)
	if _, ok := d.connections[conn]; ok {
		delete(d.connections, conn)
		close(conn.send)
		conn = nil
	}
}

func (d *Dispatcher) processInternalMessage(
	msg *message.Message,
	sender *connection,
) {
	if d.store.Length() == 0 {
		log.Print("[ERROR] Need a root auth to create new child auths")
		return
	}

	switch msg.Channel {
	case "upsub/auth/create":
		if !d.store.Find(sender.appID).CanCreate() {
			log.Print("[AUTH] Not allowed to create")
			return
		}
		d.store.CreateFromMessage(msg, sender.appID)
		break
	case "upsub/auth/update":
		if !d.store.Find(sender.appID).CanUpdate() {
			log.Print("[AUTH] Not allowed to update")
			return
		}
		d.store.UpdateFromMessage(msg)
		break
	case "upsub/auth/delete":
		if !d.store.Find(sender.appID).CanDelete() {
			log.Print("[AUTH] Not allowed to delete")
			return
		}
		d.store.DeleteFromMessage(msg)
		for conn := range d.connections {
			if conn.appID == msg.Payload {
				conn.close()
			}
		}
		break
	}
}

// ProcessMessage is parsing and routing the messages to the correct functions
func (d *Dispatcher) ProcessMessage(
	msg *message.Message,
	sender *connection,
) {
	if msg.Header == nil {
		log.Print("[INVALID MESSAGE] ", msg)
		return
	}

	msgType := msg.Type

	switch msgType {
	case message.SUBSCRIBE:
		log.Print("[SUBSCRIBE] ", msg.Payload)
		channels := strings.Split(msg.Payload, ",")
		sender.subscribe(channels)
	case message.UNSUBSCRIBE:
		log.Print("[UNSUBSCRIBE] ", msg.Payload)
		channels := strings.Split(msg.Payload, ",")
		sender.unsubscribe(channels)
	case message.PING:
		log.Print("[PING] ", msg.Payload)
		sender.send <- message.Pong()
	case message.BATCH:
		log.Print("[BATCH] ", msg.Payload)
		for _, msg := range msg.ParseBatch() {
			d.ProcessMessage(msg, sender)
		}
	case message.JSON:
		fallthrough
	case message.TEXT:
		log.Print("[RECEIVED] ", msg.Channel, " ", msg.Payload)

		if !strings.Contains(msg.Channel, ",") {
			d.Dispatch(msg, sender)
			return
		}

		channels := strings.Split(msg.Channel, ",")

		for _, channel := range channels {
			d.Dispatch(message.Text(channel, msg.Payload), sender)
		}

	}
}

// Dispatch sends messages to listening clients
func (d *Dispatcher) Dispatch(
	msg *message.Message,
	sender *connection,
) {
	if !msg.FromBroker {
		d.broker.send("upsub.dispatcher.message", msg)
	}

	if util.Contains(reservedChannels, msg.Channel) {
		d.processInternalMessage(msg, sender)
		return
	}

	for receiver := range d.connections {
		res := message.Text(msg.Channel, msg.Payload)

		if sender != nil && sender == receiver {
			continue
		}

		if sender != nil && !sender.shouldSend(receiver) {
			continue
		}

		if !receiver.shouldReceive(res) {
			continue
		}

		if sender.appID != "" {
			res.Header.Set("upsub-app-id", sender.appID)
		}

		receiver.send <- res
	}
}
