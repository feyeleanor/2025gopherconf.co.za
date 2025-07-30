package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	DialServer("tcp", "localhost:1024", func(c net.Conn) {
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
	if m, e := ReceiveMessage(c); e == nil {
		if len(m) == 0 {
			log.Println("File", n, "not found")
			return
		}
		f(m)
	}
}

func DialServer(p, a string, f func(net.Conn)) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var d net.Dialer
	if c, e := d.DialContext(ctx, p, a); e == nil {
		defer c.Close()
		f(c)
	} else {
		log.Println(e)
	}
}
