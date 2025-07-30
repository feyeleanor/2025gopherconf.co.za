package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/gob"
	"encoding/pem"
	"fmt"
	"log"
	"net"
)

type PublicKey struct {
	Label []byte
	*rsa.PublicKey
}

func main() {
	k, e := LoadPrivateKey("client_key")
	if e != nil {
		log.Fatal(e)
	}
	p := PublicKey{[]byte("served"), &k.PublicKey}

	DialUDP(":1024", func(c net.Conn) {
		defer c.Close()
		if e := SendKey(c, &p); e == nil {
			if m, e := Decrypt(k, ReadStream(c), p.Label); e == nil {
				fmt.Println(string(m))
			} else {
				log.Println(e)
			}
		}
	})
}

func DialUDP(addr string, f func(net.Conn)) {
	var a *net.UDPAddr
	var c *net.UDPConn
	var e error

	if a, e = net.ResolveUDPAddr("udp", addr); e == nil {
		if c, e = net.DialUDP("udp", nil, a); e == nil {
			defer c.Close()
			f(c)
		}
	}
	if e != nil {
		log.Println(e)
	}
}

func Decrypt(key *rsa.PrivateKey, m, l []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha1.New(), rand.Reader, key, m, l)
}

func LoadPrivateKey(f string) (r *rsa.PrivateKey, e error) {
	if block, _ := pem.Decode(LoadFile(".pem", f)); block != nil {
		if block.Type == "RSA PRIVATE KEY" {
			r, e = x509.ParsePKCS1PrivateKey(block.Bytes)
		}
	}
	return
}

func SendKey(c net.Conn, k *PublicKey) (e error) {
	var b bytes.Buffer
	if e = gob.NewEncoder(&b).Encode(k); e == nil {
		_, e = c.Write(b.Bytes())
	}
	return
}
