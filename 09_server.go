package main

import (
	"crypto/tls"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

const PROTOCOL = "tcp"
const CERT_FILE = "06_server_cert.pem"
const KEY_FILE = "06_server_key.pem"

func main() {
	dir := GetDir(os.Args[1:]...)
	AsTlsClient(CERT_FILE, KEY_FILE, func(c *tls.Config) {
		HandleTlsConnections(PROTOCOL, ADDRESS, c, func(c net.Conn) {
			CommandLoop(c, func(m string) {
				LoadFile(JSON, GetPath(dir, strings.ToLower(m)), func(j string, e error) {
					log.Printf("%v: requested file %v\n", c.RemoteAddr(), m)
					var x string
					if e == nil {
						x = strings.ReplaceAll(j, string(MESSAGE_TERMINATOR), "")
					}
					SendMessage(c, x)
				})
			})
		})
	})
}

func AsTlsClient(c, k string, f func(*tls.Config)) {
	cert, e := tls.LoadX509KeyPair(c, k)
	if e != nil {
		log.Fatal(e)
	}

	f(&tls.Config{
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

func CommandLoop(c net.Conn, f func(string)) {
	for {
		b, e := ReceiveMessage(c)
		switch e {
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
