package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

const CERT_FILE = "02_cert.pem"
const KEY_FILE = "02_key.pem"

func main() {
	dir := GetDir(os.Args[1:]...)
	cache := make(Cache)
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(cache.Read(r.URL.Path, func() ([]byte, error) {
			return LoadFile(HTML, GetPath(dir, r.URL.Path), func(h string, e error) {
				LogErrors(
					e,
					func() {
						w.Header().Set("Content-Type", "text/html")
					},
					func() {
						http.NotFound(w, r)
					})
			})
		})))
	})
	fmt.Println(http.ListenAndServeTLS(ADDRESS, CERT_FILE, KEY_FILE, nil))
}

type Cache map[string][]byte

func (c Cache) Read(k string, f func() ([]byte, error)) []byte {
	k = strings.ToLower(k)
	if _, ok := c[k]; ok {
		log.Println(k, "found in cache")
	} else {
		log.Println(k, "not stored in cache")
		if v, e := f(); e == nil {
			fmt.Printf("caching %v\n", k)
			c[k] = v
		}
	}
	return c[k]
}
