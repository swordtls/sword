package core

import (
	"crypto/x509"
	"io/ioutil"
	"log"
)

func AddTrust(pool *x509.CertPool, path string) {
	aCrt, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("ReadFile err:", err)
		return
	}
	pool.AppendCertsFromPEM(aCrt)
}
