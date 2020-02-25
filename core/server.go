package core

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sword/conf"
	"time"
)

func Sever(filename string) {
	serverConfig := conf.ParseServerConfig(filename)

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
		go ServerProxy(conn)

	}
}

func ServerProxy(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			return
		}
	}()
	defer conn.Close()
	var b = make([]byte, 1024)
	n, err := conn.Read(b[:])
	if err != nil {
		log.Println(err)
		return
	}

	if b[0] == 0x05 {
		conn.Write([]byte{0x05, 0x00})
		n, err = conn.Read(b[:])
		var host, port string
		switch b[3] {
		case 0x01:
			host = net.IPv4(b[4], b[5], b[6], b[7]).String()
		case 0x03:
			host = string(b[5 : n-2])
		case 0x04:
			host = net.IP{b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15], b[16], b[17], b[18], b[19]}.String()
		}
		port = strconv.Itoa(int(b[n-2])<<8 | int(b[n-1]))

		server, err := net.Dial("tcp", net.JoinHostPort(host, port))
		if err != nil {
			log.Println(err)
			return
		}
		defer server.Close()
		conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		go io.Copy(server, conn)
		io.Copy(conn, server)
		return
	} else {
		r := bytes.NewReader(b)
		buf := bufio.NewReader(r)
		req, err := http.ReadRequest(buf)
		if err == io.EOF {
			return
		}
		if req.Method == "CONNECT" {
			host := req.URL.Host
			server, err := net.Dial("tcp", host)
			if err != nil {
				return
			}
			conn.Write([]byte("HTTP/1.0 200 Connection Established\r\n\r\n"))
			go io.Copy(server, conn)
			io.Copy(conn, server)
			return
		} else {
			var method, host, address string
			fmt.Sscanf(string(b[:bytes.IndexByte(b[:], '\n')]), "%s%s", &method, &host)
			hostPortURL, err := url.Parse(host)
			if err != nil {
				log.Println(err)
				return
			}
			if hostPortURL.Opaque == "443" {
				address = hostPortURL.Scheme + ":443"
			} else {
				if strings.Index(hostPortURL.Host, ":") == -1 {
					address = hostPortURL.Host + ":80"
				} else {
					address = hostPortURL.Host
				}
			}
			server, err := net.Dial("tcp", address)
			if err != nil {
				log.Println(err)
				return
			}
			server.Write(b)
			go io.Copy(server, conn)
			io.Copy(conn, server)
			return

		}
	}
	log.Println(time.Now(),"unknown proxy type addr:",conn.RemoteAddr().String())
}
