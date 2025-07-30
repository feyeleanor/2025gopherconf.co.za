package main

import (
	"crypto/rand"
	"crypto/tls"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	dir := GetDir(os.Args[1:]...)
	TlsClient("server_cert.pem", "server_key.pem", func(c *tls.Config) {
		HandleTlsConnections("tcp", "localhost:1024", c, func(c net.Conn) {
			MessageLoop(c, func(m string) {
				log.Println(c.RemoteAddr(), "requested file", m)

				b := LoadFile(".json", GetPath(dir, m))
				if b != nil {
					b = DeleteAll(b, byte('\n'))
				}
				SendMessage(c, b)
			})
		})
	})
}

func TlsClient(c, k string, f func(*tls.Config)) {
	cert, e := tls.LoadX509KeyPair(c, k)
	if e != nil {
		log.Fatal(e)
	}

	f(&tls.Config{
		Rand:         rand.Reader,
		Certificates: []tls.Certificate{cert},
	})
}

func HandleTlsConnections(p, a string, c *tls.Config, f func(net.Conn)) {
	if l, e := tls.Listen(p, a, c); e == nil {
		for {
			if c, e := l.Accept(); e == nil {
				go func(c net.Conn) {
					defer c.Close()
					log.Printf("%v: connected\n", c.RemoteAddr())

					f(c)
				}(c)
			} else {
				log.Println(e)
			}
		}
	} else {
		log.Println(e)
	}
}

func MessageLoop(c net.Conn, f func(string)) {
	for {
		switch b, e := ReceiveMessage(c); e {
		case nil:
			for _, m := range Tokens(b) {
				f(m)
			}
		case io.EOF:
			log.Printf("%v: client connection dropped\n", c.RemoteAddr())
			return
		default:
			log.Printf("%v: %v\n", c.RemoteAddr(), e)
			return
		}
	}
}
