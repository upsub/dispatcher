package dispatcher

import (
	"log"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/upsub/dispatcher/src/config"
	"github.com/upsub/dispatcher/src/message"
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

	c.connection.SetReadDeadline(time.Now().Add(c.config.ReadTimeout))
	c.connection.SetPongHandler(func(string) error {
		c.connection.SetReadDeadline(time.Now().Add(c.config.ReadTimeout))
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

func (c *connection) getWildcardSubscriptions() []string {
	wildcards := []string{}

	for _, channel := range c.subscriptions {
		if strings.Contains(channel, "*") || strings.Contains(channel, ">") {
			wildcards = append(wildcards, channel)
		}
	}

	return wildcards
}

func (c *connection) shouldReceive(msg *message.Message) bool {
	channel := msg.Header.Get("upsub-channel")

	if wildcards := c.getWildcardSubscriptions(); len(wildcards) > 0 {
		channel = compareAgainstWildcardSubscriptions(wildcards, channel)
		msg.Header.Set("upsub-channel", channel)
	}

	for _, channel := range strings.Split(channel, ",") {
		if util.Contains(c.subscriptions, channel) {
			return true
		}
	}

	return false
}

func compareAgainstWildcardSubscriptions(
	wildcards []string,
	channel string,
) string {
	newChannel := channel
	for _, wildcard := range wildcards {
		if wildcardIsMatchingChannel(wildcard, channel) {
			newChannel += "," + wildcard
		}
	}

	return newChannel
}

func wildcardIsMatchingChannel(wildcard string, channel string) bool {
	wildcardParts := strings.Split(strings.Split(wildcard, ":")[0], "/")
	channelParts := strings.Split(strings.Split(channel, ":")[0], "/")

	if len(wildcardParts) > len(channelParts) {
		return false
	}

	for i, channelPart := range channelParts {
		if i > len(wildcardParts)-1 {
			// index out of bounce, not matching wildcard
			return false
		}

		if wildcardParts[i] == channelPart {
			// check if parts matches
			continue
		}

		if wildcardParts[i] == "*" {
			// check for wildcard *
			continue
		}

		if wildcardParts[i] == ">" {
			wildcardParts = createWildcardParts(wildcardParts, channelParts, i)
			continue
		}

		return false
	}

	if strings.Contains(wildcard, ":") || strings.Contains(channel, ":") {
		wildcardActions := strings.Split(wildcard, ":")
		channelActions := strings.Split(channel, ":")

		if len(wildcardActions) > 1 && len(channelActions) > 1 {
			return wildcardActions[1] == channelActions[1]
		}

		return false
	}

	return true
}

func createWildcardParts(wildcardParts []string, channelParts []string, i int) []string {
	reminders := []string{}
	newWildcardParts := []string{}

	for index, part := range wildcardParts {
		if index <= i {
			continue
		}

		reminders = append(reminders, part)
	}

	for index, part := range channelParts {
		if index > len(channelParts)-len(reminders)-1 {
			continue
		}

		newWildcardParts = append(newWildcardParts, part)
	}

	return append(newWildcardParts, reminders...)
}
