package main

import (
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	dir := GetDir(os.Args[1:]...)
	cache := make(Cache)
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if b := cache.LoadFile(".html", GetPath(dir, r.URL.Path)); b != nil {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(b)
			return
		}
		PageNotFound(w, r)
	})
	log.Println(http.ListenAndServeTLS("localhost:1024", "server_cert.pem", "server_key.pem", nil))
}

func PageNotFound(w http.ResponseWriter, r *http.Request) {
	if b := LoadFile(".html", "missing"); b != nil {
		w.Header().Set("Content-Type", "text/html")
		w.Write(b)
		http.Error(w, "", http.StatusNotFound)
	} else {
		http.NotFound(w, r)
	}
}

type Cache map[string][]byte

func (c Cache) LoadFile(t, p string) (b []byte) {
	var ok bool

	k := strings.ToLower(p)
	if b, ok = c[k]; ok {
		log.Println(k, "found in cache")
	} else {
		log.Println(k, "not found in cache")
		b = LoadFile(t, p)
		if b == nil {
			return
		}
		log.Println("caching", k)
		c[k] = b
	}
	return
}
