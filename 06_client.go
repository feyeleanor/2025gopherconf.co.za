package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	Parallelize(os.Args[1:], func(s string) {
		TlsClient("client_cert.pem", "client_key.pem", func(c *http.Client) {
			FetchWebPage(c, ServerUrl(s), func(b []byte) {
				ForEachRecord(b, func(p Person) {
					fmt.Printf("%v.json: %v [#%v] is %v\n", s, p.Name, p.Id, p.Age)
				})
			})
		})
	})
}

func FetchWebPage(c *http.Client, url string, f func([]byte)) {
	r, e := c.Get(url)
	log.Printf("fetching %v: %v\n", url, r.StatusCode)
	if e == nil {
		defer r.Body.Close()
		if r.StatusCode == http.StatusOK {
			var b []byte
			if b, e = io.ReadAll(r.Body); e == nil {
				f(b)
			}
		}
	} else {
		log.Println("fetching", url)
	}
}

func TlsClient(c, k string, f func(*http.Client)) {
	cert, _ := LoadCert(c, k)
	ca, _ := x509.SystemCertPool()
	f(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Rand:               rand.Reader,
				RootCAs:            ca,
				Certificates:       []tls.Certificate{cert},
				InsecureSkipVerify: true}}})
}

func LoadCert(c, k string) (r tls.Certificate, e error) {
	if r, e = tls.LoadX509KeyPair(c, k); e != nil {
		log.Fatalln("Error loading certificate and key file:", e)
	}
	return
}
