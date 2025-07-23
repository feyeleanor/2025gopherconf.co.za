package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
)

const CERT_FILE = "06_server_cert.pem"
const KEY_FILE = "06_server_key.pem"

func main() {
	dir := GetDir(os.Args[1:]...)

	s := NewTlsServer(ADDRESS, tls.RequestClientCert)
	//:= NewTlsServer(ADDRESS, tls.RequireAndVerifyClientCert)

	s.AddRoutes(map[string]func(http.ResponseWriter, *http.Request){
		"GET /": func(w http.ResponseWriter, r *http.Request) {
			LoadFile(JSON, GetPath(dir, r.URL.Path), func(j string, e error) {
				LogErrors(
					e,
					func() {
						w.Header().Set("Content-Type", "application/json")
						fmt.Fprint(w, j)
					},
					func() {
						http.NotFound(w, r)
					})
			})
		},
	})
	fmt.Println(s.ListenAndServeTLS(CERT_FILE, KEY_FILE))
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
			},
			Handler: http.NewServeMux(),
		}}
}
