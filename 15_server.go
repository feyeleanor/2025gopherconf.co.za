package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

func ServeFile(route, name, mime_type string) {
	b, _ := ioutil.ReadFile(name)
	http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", mime_type)
		fmt.Fprint(w, string(b))
	})
}

func main() {
	http.Handle("GET /hello", websocket.Handler(func(ws *websocket.Conn) {
		websocket.JSON.Send(ws, "Hello")
	}))

	ServeFile("GET /", "ws_hello.html", "text/html")
	ServeFile("GET /js", "ws_hello.js", "application/javascript")
	log.Println(http.ListenAndServe(":1024", nil))
}
