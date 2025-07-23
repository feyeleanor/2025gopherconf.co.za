package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const PROTOCOL = "tcp"
const CERT_FILE = "06_server_cert.pem"
const KEY_FILE = "06_server_key.pem"

func main() {
	AsTlsClient(CERT_FILE, KEY_FILE, func(c *tls.Config) {
		DialTlsServer(PROTOCOL, ADDRESS, c, func(c net.Conn) {
			for _, n := range os.Args[1:] {
				FetchFile(c, n, func(b []byte) {
					ForEachRecord(b, func(p Person) {
						fmt.Printf("%v.json: %v [#%v] is %v\n", n, p.Name, p.Id, p.Age)
					})
				})
			}
		})
	})
}

func FetchFile(c net.Conn, n string, f func([]byte)) {
	log.Println("Requesting file", n)
	SendMessage(c, n)
	if m, e := ReceiveMessage(c); e == nil {
		if len(m) == 0 {
			log.Println("File", n, "not found")
			return
		}
		f(m)
	} else {
		log.Printf("%v: %v\n", n, e)
	}
}

func AsTlsClient(c, k string, f func(*tls.Config)) {
	cert, e := tls.LoadX509KeyPair(c, k)
	if e != nil {
		log.Fatal(e)
	}

	ca, e := x509.SystemCertPool()
	if e != nil {
		log.Fatal(e)
	}

	f(&tls.Config{
		RootCAs:            ca,
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	})
}

func DialTlsServer(p, a string, c *tls.Config, f func(net.Conn)) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	d := &tls.Dialer{
		Config: c,
	}
	if c, e := d.DialContext(ctx, PROTOCOL, ADDRESS); e == nil {
		defer c.Close()

		f(c)
	} else {
		log.Fatal(e)
	}
}
