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

func main() {
	dir := GetDir(os.Args[1:]...)
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		LoadJsonFile(dir, strings.ToLower(r.URL.Path), func(j string, e error) {
			if e == nil {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, j)
			} else {
				http.NotFound(w, r)
			}
		})
	})
	fmt.Println(http.ListenAndServeTLS(ADDRESS, CERT_FILE, KEY_FILE, nil))
}

func LoadJsonFile(dir, path string, f func(string, error)) {
	j, e := ioutil.ReadFile(AddSuffix(dir + path, ".json"))
	f(string(j), e)
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
