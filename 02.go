package main
import (
	"fmt"
	"net/http"
)

const MESSAGE = "hello world"
const ADDRESS = ":1024"

const CERT_FILE = "02_cert.pem"
const KEY_FILE = "02_key.pem"

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, MESSAGE)
	})

	if e := http.ListenAndServeTLS(ADDRESS, CERT_FILE, KEY_FILE, nil); e != nil {
		fmt.Println(e)
	}
}
