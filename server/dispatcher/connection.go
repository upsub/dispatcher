package dispatcher

import (
	"log"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/upsub/dispatcher/server/config"
	"github.com/upsub/dispatcher/server/message"
	"github.com/upsub/dispatcher/server/util"
)

type connection struct {
	id            string
	name          string
	appID         string
	support       map[string]bool
	subscriptions []string
	send          chan *message.Message
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
	conn *config.Config,
	d *Dispatcher,
	support map[string]bool,
) {
	newConnection := &connection{
		id:            id,
		name:          name,
		appID:         appID,
		support:       support,
		subscriptions: []string{},
		send:          make(chan *message.Message),
		connection:    wsConnection,
		config:        conn,
		dispatcher:    d,
	}

	d.register <- newConnection
	go newConnection.read()
	go newConnection.write()
	newConnection.onConnect()
}

func (conn *connection) subscribe(channels []string) {
	for _, channel := range channels {
		conn.subscriptions = append(conn.subscriptions, channel)
	}
}

func (conn *connection) unsubscribe(channels []string) {
	var tmp []string

	for _, subscription := range conn.subscriptions {
		if !util.Contains(channels, subscription) {
			tmp = append(tmp, subscription)
		}
	}

	conn.subscriptions = tmp
}

func (conn *connection) onConnect() {
	if conn.name == "" {
		return
	}

	conn.dispatcher.Dispatch(
		message.Text("upsub/presence/"+conn.name+"/connect", ""),
		conn,
	)
}

func (conn *connection) onDisconnect() {
	if conn.name == "" {
		return
	}

	conn.dispatcher.Dispatch(
		message.Text("upsub/presence/"+conn.name+"/disconnect", ""),
		conn,
	)
}

func (conn *connection) close() {
	conn.connection.Close()
}

func (conn *connection) write() {
	ticker := time.NewTicker(conn.config.PingInterval)
	defer func() {
		ticker.Stop()
		conn.close()
	}()
	for {
		select {
		case msg, ok := <-conn.send:
			if !ok {
				conn.connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			writer, err := conn.connection.NextWriter(websocket.TextMessage)

			if err != nil {
				log.Println("[FAILED]", "Message cloudn't be written")
				return
			}

			encoded, err := msg.Encode()

			if err != nil {
				log.Print("[FAILED]", err)
				continue
			}

			writer.Write(encoded)

			if err := writer.Close(); err != nil {
				return
			}
		case <-ticker.C:
			conn.connection.SetWriteDeadline(time.Now().Add(conn.config.WriteTimeout))
			if err := conn.connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (conn *connection) read() {
	defer func() {
		defer conn.onDisconnect()
		conn.dispatcher.unregister <- conn
		conn.close()
	}()

	if conn.config.MaxMessageSize > 0 {
		conn.connection.SetReadLimit(conn.config.MaxMessageSize)
	}

	conn.connection.SetReadDeadline(time.Now().Add(conn.config.ReadTimeout))
	conn.connection.SetPongHandler(func(string) error {
		conn.connection.SetReadDeadline(time.Now().Add(conn.config.ReadTimeout))
		return nil
	})

	for {
		_, message, err := conn.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				// TODO should only print if max loglevel is set.
				log.Printf("[ERROR] %v", err)
			}
			break
		}

		conn.dispatcher.broadcast <- (func() ([]byte, *connection) {
			return message, conn
		})
	}
}

func (conn *connection) isParentToSender(sender *connection) bool {
	config := conn.config
	receiverApp := config.Apps.Find(conn.appID)
	senderApp := config.Apps.Find(sender.appID)

	if senderApp.ChildOf(receiverApp) {
		return true
	}

	return false
}

func (conn *connection) shouldSend(receiver *connection) bool {
	if conn.appID == receiver.appID {
		return true
	}

	if receiver.isParentToSender(conn) {
		return true
	}

	return false
}

func (conn *connection) getWildcardSubscriptions() []string {
	wildcards := []string{}

	for _, channel := range conn.subscriptions {
		if strings.Contains(channel, "*") || strings.Contains(channel, ">") {
			wildcards = append(wildcards, channel)
		}
	}

	return wildcards
}

func (conn *connection) shouldReceive(msg *message.Message) bool {
	channel := msg.Header.Get("upsub-channel")

	if !conn.support["wildcard"] {
		return util.Contains(conn.subscriptions, channel)
	}

	if wildcards := conn.getWildcardSubscriptions(); len(wildcards) > 0 {
		channel = compareAgainstWildcardSubscriptions(wildcards, channel)
		msg.Header.Set("upsub-channel", channel)
	}

	for _, channel := range strings.Split(channel, ",") {
		if util.Contains(conn.subscriptions, channel) {
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
