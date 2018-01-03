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
  - text
  - binary

#### Message structure
A message contains headers and a payload of events.
```json
{
  "headers": {
    "upsub-app-id": "string",
    "upsub-message-type": "string"
  },
  "payload": [{
    "event": "string",
    "channel": "string",
    "body": "string"
  }]
}
```

### Go deps
```sh
go get github.com/gorilla/websocket
go get github.com/nats-io/go-nats
```
