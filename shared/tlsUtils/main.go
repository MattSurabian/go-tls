/**
 * tlsUtils
 * This package provides a shared way to load TLS certs and keys, whether creating a
 * connection for the client or a listener for the server.
 */
package tlsUtils

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/mattsurabian/go-tls/shared/cliUtils"
	"io/ioutil"
	"net"
)

func getCipherSuites() []uint16 {
	return []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		//tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
		//tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		//tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		//tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	}
}

/**
 * loadCertificates
 * Helper method to load a specified cert and key for TLS. Both the client and the server will
 * user this method.
 */
func loadCertificates(certPath string, keyPath string) (cert tls.Certificate, certPool *x509.CertPool, err error) {
	caFile := cliUtils.GetRootCert()
	cert, err = tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return
	}

	pem, err := ioutil.ReadFile(caFile)
	if err != nil {
		return
	}

	certPool = x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pem) {
		panic("Failed appending certs")
	}

	return
}

/**
 * GetClientTLSConnected
 * Helper method called by the client to establish a connection to a remote server.
 * The connection can be used to transmit data securely.
 */
func GetClientTLSConnection() (conn *tls.Conn, err error) {
	cert, certPool, err := loadCertificates(cliUtils.GetClientTLSCertPath(), cliUtils.GetClientTLSKeyPath())
	if err != nil {
		panic(errors.New("Cannot load client TLS certs or keys, maybe run config?"))
	}

	config := &tls.Config{
		RootCAs:                certPool,
		Certificates:           []tls.Certificate{cert},
		MinVersion:             tls.VersionTLS12,
		SessionTicketsDisabled: true,
		ServerName:             cliUtils.GetRootName(),
		CipherSuites:           getCipherSuites(),
	}

	conn, err = tls.Dial("tcp", cliUtils.GetHostAndPort(), config)
	if err != nil {
		return
	}

	err = conn.Handshake()
	if err != nil {
		return
	}

	return
}

/**
 * GetServerTLSListener
 * Helper method which is called by the server so it can listen for incomming client connections.
 */
func GetServerTLSListener() (listener net.Listener) {

	cert, certPool, err := loadCertificates(cliUtils.GetServerTLSCertPath(), cliUtils.GetServerTLSKeyPath())
	if err != nil {
		panic(errors.New("Cannot load server TLS certs or keys, maybe run config?"))
	}

	config := &tls.Config{
		ClientCAs:              certPool,
		ClientAuth:             tls.RequireAndVerifyClientCert,
		Certificates:           []tls.Certificate{cert},
		MinVersion:             tls.VersionTLS12,
		SessionTicketsDisabled: true,
		CipherSuites:           []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256},
	}

	listener, err = tls.Listen("tcp", cliUtils.GetHostAndPort(), config)
	if err != nil {
		// If the server cannot open this listener, we panic, cause yo.
		panic(err)
	}

	return
}
