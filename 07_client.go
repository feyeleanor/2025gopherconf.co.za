package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const PROTOCOL = "tcp"

func main() {
	DialServer(PROTOCOL, ADDRESS, func(c net.Conn) {
		for _, n := range os.Args[1:] {
			FetchFile(c, n, func(b []byte) {
				ForEachRecord(b, func(p Person) {
					fmt.Printf("%v.json: %v [#%v] is %v\n", n, p.Name, p.Id, p.Age)
				})
			})
		}
	})
}

func FetchFile(c net.Conn, n string, f func([]byte)) {
	log.Println("Requesting file", n)
	SendMessage(c, n)
	m, e := ReceiveMessage(c)
	LogErrors(e, func() {
		if len(m) == 0 {
			log.Println("File", n, "not found")
			return
		}
		f(m)
	})
}

func DialServer(p, a string, f func(net.Conn)) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var d net.Dialer
	c, e := d.DialContext(ctx, PROTOCOL, ADDRESS)
	LogErrors(e, func() {
		defer c.Close()
		f(c)
	})
}
