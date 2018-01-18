package dispatcher

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/upsub/dispatcher/src/config"
	"github.com/upsub/dispatcher/src/util"
)

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

// CreateConnection establishes a new websocket connection
func CreateConnection(
	id string,
	appID string,
	name string,
	wsConnection *websocket.Conn,
	c *config.Config,
	d *Dispatcher,
) {
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
}

func (c *connection) subscribe(channels []string) {
	for _, channel := range channels {
		c.subscriptions = append(c.subscriptions, channel)
	}
}

func (c *connection) unsubscribe(channels []string) {
	var tmp []string

	for _, subscription := range c.subscriptions {
		if !util.Contains(channels, subscription) {
			tmp = append(tmp, subscription)
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
				log.Println("[FAILED]", "Message cloudn't be written")
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
				log.Printf("[ERROR] %v", err)
			}
			break
		}

		c.dispatcher.broadcast <- (func() ([]byte, *connection) {
			return message, c
		})
	}
}
