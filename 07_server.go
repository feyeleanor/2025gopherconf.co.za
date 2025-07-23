package main

import (
	"io"
	"log"
	"net"
	"os"
	"strings"
)

const PROTOCOL = "tcp"

func main() {
	dir := GetDir(os.Args[1:]...)
	HandleConnections(PROTOCOL, ADDRESS, func(c net.Conn) {
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
}

func HandleConnections(p, a string, f func(net.Conn)) {
	l, e := net.Listen(p, a)
	LogErrors(e, func() {
		for {
			c, e := l.Accept()
			LogErrors(e, func() {
				go func(c net.Conn) {
					defer c.Close()
					log.Printf("%v: connected\n", c.RemoteAddr())

					f(c)
				}(c)
			})
		}
	})
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
