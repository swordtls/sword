package core

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"log"
	"net"
	"sword/conf"
)

func ProxyServer(filename string) {
	serverConfig := conf.ParseProxyServerConfig(filename)

	cert, err := tls.LoadX509KeyPair(serverConfig.ServerPerm, serverConfig.ServerKey)
	if err != nil {
		log.Println(err)
		return
	}

	certBytes, err := ioutil.ReadFile(serverConfig.ClientPerm)
	if err != nil {
		panic("Unable to read cert.pem")
	}
	clientCertPool := x509.NewCertPool()
	AddTrust(clientCertPool, serverConfig.CaCert)
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		panic("failed to parse root certificate")
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCertPool,
	}
	l, err := tls.Listen("tcp", serverConfig.ListenPort, config)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go DoClientProxyNoTls(conn, serverConfig.RedirectAddress)

	}
}

func DoClientProxyNoTls(conn net.Conn, serverAddress string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			return
		}
	}()
	defer conn.Close()
	var err error
	var server net.Conn
	server, err = net.Dial("tcp", serverAddress)
	if err != nil {
		log.Println(err)
		return
	}
	defer server.Close()
	go io.Copy(server, conn)
	io.Copy(conn, server)
}
