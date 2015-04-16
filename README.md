# go-tls
Minimal implementation of a configurable client/server using a tls tunnel.

## Deps
To play with this repository you must have `Go` [installed on your system](https://golang.org/doc/install).

While this project plays nicely with all of Go's built in tooling it also provides a vendor script and Makefile
so it doesn't matter if you're using a single global `$GOPATH` or not.

## Building
To play with this repository using Go's built in tooling it should be cloned to the "expected"
location in your GoPath, the easiest way to do this is with `go get github.com/mattsurabian/go-tls`.

You can then run `go get` and `go build` in the `client` and `server` directories.

You can also just check it out to any location you like and use `make` to build the client and
server binaries in their respective folders.

## Cipher Suites
A list of NIST "should" ciphers is provided but since the entirety of the client/server relationship is
represented it's not necessary to support more than one cipher suite. The default choice is `tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256`
as I believe it to be the most secure suite available for TLS right now.

## Minting Certs and Keys
I used [@bnagy's Enough repo](https://github.com/bnagy/enough) and the included `tlspark` tool to create
the test certs included in this repo. The `root-subject` configuration flag corresponds to the `name` flag
passed into `tlspark`.

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