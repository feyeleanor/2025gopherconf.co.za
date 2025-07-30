package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	dir := GetDir(os.Args[1:]...)
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if b, e := ioutil.ReadFile(AddSuffix(GetPath(dir, r.URL.Path), ".html")); e == nil {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(b)
		} else {
			log.Println(e)
			http.NotFound(w, r)
		}
	})
	log.Println(http.ListenAndServeTLS("localhost:1024", "server_cert.pem", "server_key.pem", nil))
}
