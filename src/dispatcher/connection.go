package dispatcher

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/upsub/dispatcher/src/config"
	"github.com/upsub/dispatcher/src/util"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type connection struct {
	id            string
	name          string
	appID         string
	subscriptions []string
	send          chan []byte
	connection    *websocket.Conn
	config        *config.Config
	dispatcher    *Dispatcher
}

func createConnection(
	id string,
	appID string,
	name string,
	wsConnection *websocket.Conn,
	c *config.Config,
	d *Dispatcher,
) *connection {
	newConnection := &connection{
		id:            id,
		name:          name,
		appID:         appID,
		subscriptions: []string{},
		send:          make(chan []byte, 256),
		connection:    wsConnection,
		config:        c,
		dispatcher:    d,
	}

	d.register <- newConnection
	go newConnection.read()
	go newConnection.write()

	return newConnection
}

// UpgradeHandler http request to websocket protocol
func (d *Dispatcher) UpgradeHandler(
	config *config.Config,
	dispatcher *Dispatcher,
	w http.ResponseWriter,
	r *http.Request,
) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		return
	}

	createConnection(
		r.Header.Get("Sec-Websocket-Key"),
		r.Header.Get("upsub-app-id"),
		r.Header.Get("upsub-connection-name"),
		conn,
		config,
		dispatcher,
	)
}

func (c *connection) subscribe(channel string) {
	c.subscriptions = append(c.subscriptions, channel)
}

func (c *connection) unsubscribe(channel string) {
	var tmp []string

	for _, current := range c.subscriptions {
		if current != channel {
			tmp = append(tmp, channel)
		}
	}

	c.subscriptions = tmp
}

func (c *connection) hasSubscription(channel string) bool {
	return util.Contains(c.subscriptions, channel)
}

func (c *connection) write() {
	ticker := time.NewTicker(c.config.PingInterval)
	defer func() {
		ticker.Stop()
		c.connection.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			writer, err := c.connection.NextWriter(websocket.TextMessage)

			if err != nil {
				log.Println("[Faild]", "Message cloudn't be written")
				return
			}

			writer.Write(message)

			if err := writer.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.connection.SetWriteDeadline(time.Now().Add(c.config.WriteTimeout))
			if err := c.connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *connection) read() {
	defer func() {
		c.dispatcher.unregister <- c
		c.connection.Close()
	}()

	if c.config.MaxMessageSize > 0 {
		c.connection.SetReadLimit(c.config.MaxMessageSize)
	}

	c.connection.SetReadDeadline(time.Now().Add(c.config.Timeout))
	c.connection.SetPongHandler(func(string) error {
		c.connection.SetReadDeadline(time.Now().Add(c.config.Timeout))
		return nil
	})

	for {
		_, message, err := c.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				// TODO should only print if max loglevel is set.
				log.Printf("[Error] %v", err)
			}
			break
		}

		c.dispatcher.broadcast <- (func() ([]byte, *connection) {
			return message, c
		})
	}
}
