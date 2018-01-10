package dispatcher

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type client struct {
	id            string
	name          string
	appID         string
	subscriptions []string
	send          chan []byte
	connection    *websocket.Conn
	config        *config
	dispatcher    *dispatcher
}

func createClient(
	id string,
	appID string,
	name string,
	connection *websocket.Conn,
	c *config,
	d *dispatcher,
) *client {
	client := &client{
		id:            id,
		name:          name,
		appID:         appID,
		subscriptions: []string{},
		send:          make(chan []byte, 256),
		connection:    connection,
		config:        c,
		dispatcher:    d,
	}

	d.register <- client
	go client.read()
	go client.write()

	return client
}

func upgradeHandler(
	config *config,
	dispatcher *dispatcher,
	w http.ResponseWriter,
	r *http.Request,
) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		return
	}

	createClient(
		r.Header.Get("Sec-Websocket-Key"),
		r.Header.Get("upsub-app-id"),
		r.Header.Get("upsub-connection-name"),
		conn,
		config,
		dispatcher,
	)
}

func (c *client) subscribe(channel string) {
	c.subscriptions = append(c.subscriptions, channel)
}

func (c *client) unsubscribe(channel string) {
	var tmp []string

	for _, current := range c.subscriptions {
		if current != channel {
			tmp = append(tmp, channel)
		}
	}

	c.subscriptions = tmp
}

func (c *client) write() {
	ticker := time.NewTicker(c.config.pingInterval)
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
			c.connection.SetWriteDeadline(time.Now().Add(c.config.writeTimeout))
			if err := c.connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *client) read() {
	defer func() {
		c.dispatcher.unregister <- c
		c.connection.Close()
	}()

	if c.config.maxMessageSize > 0 {
		c.connection.SetReadLimit(c.config.maxMessageSize)
	}

	c.connection.SetReadDeadline(time.Now().Add(c.config.timeout))
	c.connection.SetPongHandler(func(string) error {
		c.connection.SetReadDeadline(time.Now().Add(c.config.timeout))
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

		c.dispatcher.broadcast <- (func() ([]byte, *client) {
			return message, c
		})
	}
}
