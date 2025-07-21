package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const ADDRESS = ":1024"
const PROTOCOL = "tcp"
const CERT_FILE = "06_server_cert.pem"
const KEY_FILE = "06_server_key.pem"

const MESSAGE_TERMINATOR = '\n'
const MESSAGE_WHITESPACE = ' '

type Person struct {
	Id		int
	Name	string
	Age		string
}

func main() {
	PrepareTls(CERT_FILE, KEY_FILE, func(c *tls.Config) {
		DialTlsServer(PROTOCOL, ADDRESS, c, func(c net.Conn) {
			for _, n := range os.Args[1:] {
				FetchFile(c, n, func(s string) {
					ForEachRecord(s, func(p Person) {
						fmt.Printf("%v.json: %v [#%v] is %v\n", n, p.Name, p.Id, p.Age)
					})
				})
			}
		})
	})
}

func ForEachRecord(s string, f func(Person)) {
	var e error
	r := []Person{}
	if e = json.Unmarshal([]byte(s), &r); e == nil {
		for _, v := range r {
			f(v)
		}
	}
}

func FetchFile(c net.Conn, n string, f func(string)) {
	log.Println("Requesting file", n)
	SendMessage(c, n)
	if m, e := ReceiveMessage(c); e == nil {
		if len(m) == 0 {
			log.Println("File", n, "not found")
			return
		}
		f(m)
	} else {
		log.Printf("%v: %v\n", n, e)
	}
}

func PrepareTls(c, k string, f func(*tls.Config)) {
	cert, e := tls.LoadX509KeyPair(c, k)
	if e != nil {
		log.Fatal(e)
	}

	f(&tls.Config {
		Certificates: []tls.Certificate{ cert },
		InsecureSkipVerify: true,
	})
}

func DialTlsServer(p, a string, c *tls.Config, f func(net.Conn)) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 30)
	defer cancel()

	d := &tls.Dialer {
		Config: c,
	}
	if c, e := d.DialContext(ctx, PROTOCOL, ADDRESS); e == nil {
		defer c.Close()

		f(c)
	} else {
		log.Fatal(e)
	}
}

func ReceiveMessage(c net.Conn) (s string, e error) {
	if s, e = bufio.NewReader(c).ReadString(MESSAGE_TERMINATOR); e == nil {
		s = s[:len(s) - 1]
	}
	return
}

func SendMessage(c net.Conn, s ...any) {
	for _, v := range s {
		switch v := v.(type) {
		case []byte:
			c.Write(v)
			c.Write([]byte(string(MESSAGE_TERMINATOR)))
		case string:
			SendMessage(c, []byte(v))
		case rune:
			SendMessage(c, string(v))
		case fmt.Stringer:
			SendMessage(c, v.String())
		default:
			log.Printf("unable to send message [%T] %v\n", v, v)
		}
	}
}