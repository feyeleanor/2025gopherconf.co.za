package main

import (
	"fmt"

	"golang.org/x/net/websocket"
)

func main() {
	if ws, e := websocket.Dial("ws://localhost:1024/hello", "", "https://localhost/"); e == nil {
		var s string
		if e := websocket.JSON.Receive(ws, &s); e == nil {
			fmt.Println(s)
		}
	}
}
