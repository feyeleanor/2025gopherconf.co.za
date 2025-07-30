package main

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	TlsClient("client_cert.pem", "client_key.pem", func(c *tls.Config) {
		DialTlsServer("tcp", "localhost:1024", c, func(c net.Conn) {
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
	}
}

func TlsClient(c, k string, f func(*tls.Config)) {
	cert, e := tls.LoadX509KeyPair(c, k)
	if e != nil {
		log.Fatal(e)
	}

	f(&tls.Config{
		Rand:               rand.Reader,
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
	if c, e := d.DialContext(ctx, p, a); e == nil {
		defer c.Close()

		f(c)
	} else {
		log.Fatal(e)
	}
}
