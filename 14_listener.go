package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/gob"
	"log"
	"net"
)

type PublicKey struct {
	Label []byte
	*rsa.PublicKey
}

func main() {
	ListenForUDP(":1024", func(c *net.UDPConn) {
		RespondToUDP(c, func(a *net.UDPAddr, b []byte) {
			var k PublicKey
			if e := gob.NewDecoder(bytes.NewBuffer(b)).Decode(&k); e == nil {

				if m, e := Encrypt(&k, []byte("Hello World")); e == nil {
					c.WriteToUDP(m, a)
				}
			}
			return
		})

	})
}

func RespondToUDP(c *net.UDPConn, f func(*net.UDPAddr, []byte)) {
	for b := make([]byte, 1024); ; b = make([]byte, 1024) {
		if _, a, e := c.ReadFromUDP(b); e == nil {
			go f(a, b)
		} else {
			log.Println(e)
		}
	}
}

func ListenForUDP(addr string, f func(*net.UDPConn)) {
	var a *net.UDPAddr
	var c *net.UDPConn
	var e error

	if a, e = net.ResolveUDPAddr("udp", addr); e == nil {
		if c, e = net.ListenUDP("udp", a); e == nil {
			f(c)
		}
	}
	if e != nil {
		log.Println(e)
	}
}

func Encrypt(k *PublicKey, m []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha1.New(), rand.Reader, k.PublicKey, m, k.Label)
}
