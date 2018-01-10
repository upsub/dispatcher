# upsub

## Why?

## Authentication

- App id
- Secret key
- Public key
- Authenticate with origin from browser
- Protocol JSON web tokens?

## Message specification

Message types:
  - subscribe
  - unsubscribe
  - ping
  - pong
  - batch
  - text
  - binary

#### Message structure
A message contains headers and a payload.
```json
{
  "headers": {
    "upsub-message-type": "string",
    "upsub-channel":      "string|optional",
  },
  "payload": "string"
}
```


Example of upsub messages
```js
// Subscribe message
{
  "headers": {
    "upsub-message-type": "subscribe"
  },
  "payload": "some-channel"
}

// Text message
{
  "headers": {
    "upsub-message-type": "text",
    "upsub-channel": "/hello"
  },
  "payload": "Hello world!"
}

// Batch message
{
  "headers": {
    "upsub-message-type": "batch",
  },
  "payload": JSON.stringify([{
    "headers": {
      "upsub-message-type": "text",
      "upsub-channel": "/hello"
    },
    "payload": "Hello"
  }, {
    "headers": {
      "upsub-message-type": "text",
      "upsub-channel": "/world"
    },
    "payload": "World"
  }])
}
```

### Go deps
```sh
go get github.com/gorilla/websocket
go get github.com/nats-io/go-nats
```
