package main

import (
	"fmt"
	"net/http"
	"os"
)

const CERT_FILE = "02_cert.pem"
const KEY_FILE = "02_key.pem"

func main() {
	dir := GetDir(os.Args[1:]...)
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		LoadFile(HTML, GetPath(dir, r.URL.Path), func(h string, e error) {
			LogErrors(
				e,
				func() {
					w.Header().Set("Content-Type", "text/html")
					fmt.Fprint(w, h)
				},
				func() {
					http.NotFound(w, r)
				})
		})
	})
	fmt.Println(http.ListenAndServeTLS(ADDRESS, CERT_FILE, KEY_FILE, nil))
}
