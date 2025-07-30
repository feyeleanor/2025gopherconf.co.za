package main

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"log"

	"golang.org/x/net/websocket"
)

func main() {
	DialWebSocket("wss://localhost:1024/hello", "https://localhost:1024", func(ws *websocket.Conn) {
		var s string
		if e := websocket.JSON.Receive(ws, &s); e == nil {
			fmt.Println(s)
		} else {
			log.Println(e)
		}

	})
}

func DialWebSocket(url, o string, f func(*websocket.Conn)) {
	var c *websocket.Config
	var e error
	var ws *websocket.Conn

	if c, e = websocket.NewConfig(url, o); e == nil {
		c.TlsConfig = &tls.Config{
			Rand:               rand.Reader,
			InsecureSkipVerify: true}

		if ws, e = websocket.DialConfig(c); e == nil {
			f(ws)
		}
	}
	if e != nil {
		log.Println(e)
	}
}
