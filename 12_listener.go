package main

import (
	"log"
	"net"
	"syscall"
)

func main() {
	ListenForUDP(":1024", func(c *net.UDPConn) {
		RespondToUDP(c, func(a *net.UDPAddr, _ []byte) {
			c.WriteToUDP([]byte("Hello World\n"), a)
		})
	})
}

func RespondToUDP(c *net.UDPConn, f func(*net.UDPAddr, []byte)) {
	log.Println("receive buffer size", GetSocketOption(c, syscall.SO_RCVBUF))
	for b := make([]byte, 1024); ; b = make([]byte, 1024) {
		log.Println("waiting to receive message...")
		if _, a, e := c.ReadFromUDP(b); e == nil {
			log.Println("received message", string(b))
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
