FROM golang:latest
WORKDIR /go/src/github.com/upsub/dispatcher
ADD . .
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure
EXPOSE 4400
CMD ["go", "run", "main.go"]
