package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

const URL = "https://localhost:1024/"

type Person struct {
	Id		int `json:"id"`
	Name	string `json:"name"`
	Age		string
}

func main() {
	var w sync.WaitGroup

	c := &http.Client {
		Transport: 	&http.Transport {
			TLSClientConfig: &tls.Config {
				InsecureSkipVerify: true }}}

/*
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
*/

	for _, n := range os.Args[1:] {
		w.Add(1)
		go func(n string) {
			defer w.Done()

			FetchWebPage(c, URL + n, func(b []byte) {
				ForEachRecord(b, func(p Person) {
					fmt.Printf("%v.json: %v [#%v] is %v\n", n, p.Name, p.Id, p.Age)
				})
			})
		}(n)
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
	log.Printf("fetching %v: %v\n", url, r.StatusCode)
	if e == nil {
		defer r.Body.Close()
		if r.StatusCode == http.StatusOK {
			var b []byte
			if b, e = io.ReadAll(r.Body); e == nil {
				f(b)
			}
		}
	}
}
