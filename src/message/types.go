package message

const (
	// TextMessage normal upsub message
	TextMessage = "text"

	// SubscripeMessage subscribe to topic
	SubscripeMessage = "subscribe"

	// UnsubscribeMessage remove topic subscription
	UnsubscribeMessage = "unsubscribe"

	// PingMessage ping the upsub service, should return a pong message
	PingMessage = "ping"

	// PongMessage pong message is return when a ping message is received
	PongMessage = "pong"
)
