package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const ADDRESS = ":1024"
const CERT_FILE = "02_cert.pem"
const KEY_FILE = "02_key.pem"

const ROOT_DIR = "."
const BACKSLASH = "/"

type Cache map[string] []byte

func (c Cache) Read(k string, f func() ([]byte, error)) []byte {
	if _, ok := c[k]; !ok {
		fmt.Printf("%v not stored in cache\n", k)
		if v, e := f(); e == nil {
			fmt.Printf("caching %v\n", k)
			c[k] = v
		}
	}
	return c[k]
}

func main() {
	dir := GetDir(os.Args[1:]...)
	cache := make(Cache)
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		path := strings.ToLower(r.URL.Path)
		fmt.Fprint(w, string(cache.Read(path, func() ([]byte, error) {
			return LoadHtmlFile(dir, path, func(h string, e error) {
				if e == nil {
					w.Header().Set("Content-Type", "text/html")
				} else {
					http.NotFound(w, r)
				}
			})
		})))
	})
	fmt.Println(http.ListenAndServeTLS(ADDRESS, CERT_FILE, KEY_FILE, nil))
}

func LoadHtmlFile(dir, path string, f func(string, error)) (h []byte, e error) {
	h, e = ioutil.ReadFile(AddSuffix(dir + path, ".html"))
	f(string(h), e)
	return h, e
}

func GetDir(s ...string) (d string) {
	if len(s) > 0 {
		d = RemoveDuplicates(s[0], BACKSLASH)
	} else {
		d = ROOT_DIR
	}
	return AddSuffix(d, BACKSLASH)
}

func AddSuffix(v, s string) string {
	if !strings.HasSuffix(v, s) {
		return v + s
	}
	return v
}

func RemoveDuplicates(s, sep string) string {
	return strings.Join(strings.Split(s, sep), sep)
}
