# go-tls
Minimal implementation of a configurable client/server using a tls tunnel. This repo is for reference
and to explore some [interesting cipher suite issues](https://github.com/MattSurabian/go-tls/blob/master/shared/tlsUtils/main.go#L62) I noticed while implementing TLS on another project.

##Deps
To play with this repository you must have `Go` [installed on your system](https://golang.org/doc/install).

While this project plays nicely with all of Go's built in tooling it also provides a vendor script and Makefile
so it doesn't matter if you're using a single global `$GOPATH` or not.

##Building
To play with this repository using Go's built in tooling it should be cloned to the "expected"
location in your GoPath, the easiest way to do this is with `go get github.com/mattsurabian/go-tls`.

You can then run `go get` and `go build` in the `client` and `server` directories.

You can also just check it out to any location you like and use `make` to build the client and
server binaries in their respective folders.

## Client

The client supports two commands: `config` and `send`.

### config
The config command prompts the user for several values necessary to establish a TLS tunnel to the
server. Specifically: the address of the server, the port the server is listening on, a root cert,
a client TLS cert and the corresponding key.

### send
The send command expects a string to send to the server: `./client send "some message"`

## Server

The server supports two commands: `config` and `start`

### config
The config command prompts the user for several values necessary to start listening for incoming
TLS connections from clients. Specifically: the port the server should listen on, the root cert,
a server TLS cert and the corresponding key.

### start
The start command opens a port and starts listening for incoming connections from clients: `./server start`.
Any messages it receives will be logged to `STDOUT`.