package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"log"
	"net"
)

const AES_KEY = "0123456789012345"

func main() {
	DialUDP(":1024", func(c net.Conn) {
		SendMessage(c, NewIV())
		RequestMessage(c, func(b []byte) {
			if m, e := Decrypt(b, AES_KEY); e == nil {
				fmt.Println(string(m))
			} else {
				log.Println(e)
			}
		})
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

func RequestMessage(c net.Conn, f func([]byte)) {
	if _, e := conn.Write([]byte("\n")); e == nil {
		if m := ReadStream(c); m != nil {
			f(m)
		}
	} else {
		log.Println(e)
	}
	return
}

func NewIV() (b []byte) {
	b = make([]byte, aes.BlockSize)
	if _, e := rand.Read(b); e != nil {
		panic(e)
	}
	return
}

func Decrypt(m []byte, k string) (r []byte, e error) {
	var b cipher.Block
	if b, e = aes.NewCipher([]byte(k)); e == nil {
		var iv []byte
		iv, m = Unpack(m)
		c := cipher.NewCBCDecrypter(b, iv)
		r = make([]byte, len(m))
		c.CryptBlocks(r, m)
	}
	return
}

func Unpack(m []byte) (iv, r []byte) {
	return m[:aes.BlockSize], m[aes.BlockSize:]
}
