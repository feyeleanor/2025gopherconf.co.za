package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

const URL = "https://localhost:1024/"
const CERT_FILE = "06_client_cert.pem"
const KEY_FILE = "06_client_key.pem"

type Person struct {
	Id		int `json:"id"`
	Name	string `json:"name"`
	Age		string
}

func main() {
	var w sync.WaitGroup

	for _, f := range os.Args[1:] {
		w.Add(1)
		go func(f string) {
			defer w.Done()

			PrepareTls(func(c *http.Client) {
				FetchWebPage(c, URL + f, func(b []byte) {
					ForEachRecord(b, func(p Person) {
						fmt.Printf("%v.json: %v [#%v] is %v\n", f, p.Name, p.Id, p.Age)
					})
				})
			})
		}(f)
	}
	w.Wait()
}

func ForEachRecord(b []byte, f func(Person)) {
	var e error
	r := []Person{}
	if e = json.Unmarshal(b, &r); e == nil {
		for _, v := range r {
			f(v)
		}
	}
}

func FetchWebPage(c *http.Client, url string, f func([]byte)) {
	r, e := c.Get(url)
	if e == nil {
		log.Printf("fetching %v: %v\n", url, r.StatusCode)
		defer r.Body.Close()
		if r.StatusCode == http.StatusOK {
			var b []byte
			if b, e = io.ReadAll(r.Body); e == nil {
				f(b)
			}
		}
	} else {
		log.Printf("fetching %v: %v\n", url, e)
	}
}

func PrepareTls(f func(*http.Client)) {
	cert, _ := LoadCert(CERT_FILE, KEY_FILE)
	ca, _ := x509.SystemCertPool()
	f(&http.Client {
		Transport: 	&http.Transport {
			TLSClientConfig: &tls.Config {
				RootCAs: ca,
				Certificates: []tls.Certificate{ cert },
				InsecureSkipVerify: true }}})
}

func LoadCert(c, k string) (r tls.Certificate, e error) {
	r, e = tls.LoadX509KeyPair(c, k)
	if e != nil {
		log.Fatalf("Error loading certificate and key file: %v\n", e)
	}
	return
}
