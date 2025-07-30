package main

import "net"

func main() {
	if a, e := net.ResolveUDPAddr("udp", ":1024"); e == nil {
		if c, e := net.ListenUDP("udp", a); e == nil {
			for b := make([]byte, 1024); ; b = make([]byte, 1024) {
				if _, cl, e := c.ReadFromUDP(b); e == nil {
					go func(c *net.UDPConn, a *net.UDPAddr) {
						c.WriteToUDP([]byte("Hello World\n"), a)
					}(c, cl)
				}
			}
		}
	}
}
