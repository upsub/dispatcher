# Dispatcher

> A high performance Pub/Sub messaging server for the Web and Cloud.

UpSub's Dispatcher is a bidirectional commnunication server based on
the WebSocket protocol.

## Getting Started
Pull and run the latest version of the Dispatcher from docker hub with the
command below.
```sh
docker run -d -p 4400:4400 upsub/dispatcher
```
> The command starts a dispatcher instance which listens on `localhost:4400`.

#### Docker Compose
The best way to get the Dispatcher integrated in your dev environment is to use
[`docker-compose`](https://docs.docker.com/compose/overview/). Down below is
an example of a simple `docker-compose.yml` that configures a dispatcher
instance.
```yml
version: '3'
services:
  upsub/dispatcher:
    image: upsub/dispatcher
    ports:
      - '4400:4400'
```

## Configuration
The Dispatcher is configurable through environment variables, all available
configuration options will be listed below.

- `PORT (default: 4400)`: The port the dispatcher is listing for messages.
- `MAX_MESSAGE_SIZE (default: no limit)`: The maximum size of a message in bytes.
- `CONNECTION_TIMEOUT (default: 10s)`: The dispatcher will reject a client if it exceeds the connection timeout.
- `READ_TIMEOUT (default: 10s)`: The dispatcher will terminate message if it wasn't received within the timeout.
- `WRITE_TIMEOUT (default: 10s)`: The dispatcher will terminate message if it couldn't be written within the timeout.

### Authentication configuration
The dispatcher can authenticate connections with an `App ID`, `Secret Key` if it's normal client connection.
If it's a client from a browser environment, the authentication method will be `App ID`, `Public Key`
and the `origin` of the request.

- `AUTH_APP_ID`: Should be a string which identifies your application.
- `AUTH_SECRET`: Should be a sha256 key.
- `AUTH_PUBLIC`: Should be a sha256 key.
- `AUTH_ORIGINS`: A comma separated string of the allowed origins


## License
Released under the [MIT license](https://github.com/upsub/dispatcher/blob/master/LICENSE)
