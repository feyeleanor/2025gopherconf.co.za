package main

import (
	"crypto/rand"
	"crypto/tls"
	"log"
	"net/http"
	"os"
)

func main() {
	dir := GetDir(os.Args[1:]...)

	s := NewTlsServer("localhost:1024", tls.RequestClientCert)
	//:= NewTlsServer("localhost:1024", tls.RequireAndVerifyClientCert)

	s.AddRoutes(map[string]func(http.ResponseWriter, *http.Request){
		"GET /hello": func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.RemoteAddr, "requested /hello")

			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.Write([]byte("hello world"))
		},

		"GET /": func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.RemoteAddr, "requested file", r.URL.Path)

			b := LoadFile(".json", GetPath(dir, r.URL.Path))
			if b == nil {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(b)
		},
	})
	log.Println(s.ListenAndServeTLS("server_cert.pem", "server_key.pem"))
}

type TlsServer struct {
	http.Server
}

func (s TlsServer) AddRoutes(m map[string]func(http.ResponseWriter, *http.Request)) {
	for r, f := range m {
		s.Handler.(*http.ServeMux).HandleFunc(r, f)
	}
}

func NewTlsServer(addr string, auth tls.ClientAuthType) *TlsServer {
	return &TlsServer{
		http.Server{
			Addr: addr,
			TLSConfig: &tls.Config{
				ClientAuth: auth,
				Rand:       rand.Reader,
			},
			Handler: http.NewServeMux(),
		}}
}
