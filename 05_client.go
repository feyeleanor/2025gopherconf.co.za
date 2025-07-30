package main

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	Parallelize(os.Args[1:], func(s string) {
		TlsClient(func(c *http.Client) {
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
			if b, e := io.ReadAll(r.Body); e == nil {
				f(b)
			} else {
				log.Println(e)
			}
		}
	} else {
		log.Println(e)
	}
}

func TlsClient(f func(*http.Client)) {
	f(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Rand:               rand.Reader,
				InsecureSkipVerify: true}}})
}

/*
	There is an alternative approach where we turn off Cert chain verification for
	the default http.ServeMux

		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}

	However this affects all connections using the default ServeMux which may lead
	to unexpected results
*/
