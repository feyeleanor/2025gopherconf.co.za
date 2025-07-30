package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	dir := GetDir(os.Args[1:]...)
	HandleConnections("tcp", "localhost:1024", func(c net.Conn) {
		MessageLoop(c, func(m string) {
			log.Println(c.RemoteAddr(), "requested file", m)

			b := LoadFile(".json", GetPath(dir, m))
			if b != nil {
				b = DeleteAll(b, byte('\n'))
			}
			SendMessage(c, b)
		})
	})
}

func HandleConnections(p, a string, f func(net.Conn)) {
	if l, e := net.Listen(p, a); e == nil {
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
