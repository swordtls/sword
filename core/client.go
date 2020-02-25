package core

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"log"
	"net"
	"sword/conf"
)

var tlsConf *tls.Config

func Client(filename string) {
	clientConfig := conf.ParseClientConfig(filename)

	cert, err := tls.LoadX509KeyPair(clientConfig.ClientPerm, clientConfig.ClientKey)
	if err != nil {
		log.Println(err)
		return
	}
	clientCertPool := x509.NewCertPool()
	AddTrust(clientCertPool, clientConfig.CaCert)
	tlsConf = &tls.Config{
		RootCAs:            clientCertPool,
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	l, err := net.Listen("tcp", clientConfig.ListenPort)
	if err != nil {
		log.Panic(err)
	}
	for {
		client, err := l.Accept()
		if err != nil {
			log.Println(err)
		}
		go DoClientProxy(client, clientConfig.ServerAddress)
	}
}
func DoClientProxy(conn net.Conn, serverAddress string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			return
		}
	}()
	defer conn.Close()
	var err error
	var server net.Conn
	server, err = tls.Dial("tcp", serverAddress, tlsConf)
	if err != nil {
		log.Println(err)
		return
	}
	defer server.Close()
	go io.Copy(server, conn)
	io.Copy(conn, server)
}
