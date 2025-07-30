package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", Hello)
	if e := http.ListenAndServe("localhost:1024", nil); e != nil {
		fmt.Println(e)
	}
}

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "hello world")
}
