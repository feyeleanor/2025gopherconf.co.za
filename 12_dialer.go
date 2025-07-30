package main

import (
	"fmt"
	"log"
	"net"
	"syscall"
)

func main() {
	DialUDP(":1024", func(c net.Conn) {
		log.Println("send packet to peer")
		SendMessage(c, "")
		if m, e := ReceiveMessage(c); e == nil {
			log.Println("received message", m)
			fmt.Println(string(m))
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
			log.Println("send buffer size", GetSocketOption(c, syscall.SO_SNDBUF))
			f(c)
		}
	}
	if e != nil {
		log.Println(e)
	}
}
