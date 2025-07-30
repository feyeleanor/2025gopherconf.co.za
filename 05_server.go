package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	dir := GetDir(os.Args[1:]...)
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr, "requested file", r.URL.Path)

		b := LoadFile(".json", GetPath(dir, r.URL.Path))
		if b == nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	})
	log.Println(http.ListenAndServeTLS("localhost:1024", "server_cert.pem", "server_key.pem", nil))
}
