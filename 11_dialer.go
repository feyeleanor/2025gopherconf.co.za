package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	var a *net.UDPAddr
	var c *net.UDPConn
	var e error
	var m string

	if a, e = net.ResolveUDPAddr("udp", ":1024"); e == nil {
		if c, e = net.DialUDP("udp", nil, a); e == nil {
			defer c.Close()

			if _, e = c.Write([]byte{}); e == nil {
				if m, e = bufio.NewReader(c).ReadString('\n'); e == nil {
					fmt.Print(m)
				}
			}
		}
	}
	if e != nil {
		log.Println(e)
	}

}
