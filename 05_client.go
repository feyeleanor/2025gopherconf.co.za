package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	Parallelize(os.Args[1:], func(s string) {
		AsTlsClient(func(c *http.Client) {
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
	LogErrors(
		e,
		func() {
			defer r.Body.Close()
			if r.StatusCode == http.StatusOK {
				b, e := io.ReadAll(r.Body)
				LogErrors(e, func() {
					f(b)
				})
			}
		})
}

/*
	There is an alternative approach where we turn off Cert chain verification for
	the default http.ServeMUX

		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
*/

func AsTlsClient(f func(*http.Client)) {
	f(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true}}})
}
