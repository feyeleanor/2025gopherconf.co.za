package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("hello world"))
	})

	if e := http.ListenAndServeTLS("localhost:1024", "server_cert.pem", "server_key.pem", nil); e != nil {
		log.Println(e)
	}
}
