FROM golang:latest
WORKDIR /go/src/github.com/upsub/dispatcher
ADD . .
RUN go get \
  github.com/gorilla/websocket \
  github.com/nats-io/go-nats
EXPOSE 5000
CMD ["go", "run", "main.go"]
