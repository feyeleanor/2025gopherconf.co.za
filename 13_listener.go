package main

import (
	"crypto/aes"
	"crypto/cipher"
	"log"
	"net"
)

const AES_KEY = "0123456789012345"

func main() {
	ListenForUDP(":1024", func(c *net.UDPConn) {
		RespondToUDP(c, func(a *net.UDPAddr, iv []byte) {
			if m, e := Encrypt("Hello World", AES_KEY, iv); e == nil {
				c.WriteToUDP(m, a)
			}
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

func Encrypt(m, k string, iv []byte) (o []byte, e error) {
	if o, e = PaddedBuffer([]byte(m)); e == nil {
		var b cipher.Block
		if b, e = aes.NewCipher([]byte(k)); e == nil {
			o = CryptBlocks(o, iv, b)
		}
	}
	return
}

func PaddedBuffer(m []byte) (b []byte, e error) {
	b = append(b, m...)
	if p := len(b) % aes.BlockSize; p != 0 {
		p = aes.BlockSize - p
		b = append(b, make([]byte, p)...) // padding with NUL!!!!
	}
	return
}

func CryptBlocks(b, iv []byte, c cipher.Block) (o []byte) {
	o = make([]byte, aes.BlockSize+len(b))
	copy(o, iv)
	enc := cipher.NewCBCEncrypter(c, o[:aes.BlockSize])
	enc.CryptBlocks(o[aes.BlockSize:], b)
	return
}
