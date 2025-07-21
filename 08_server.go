package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

const ADDRESS = ":1024"
const PROTOCOL = "tcp"
const CERT_FILE = "06_server_cert.pem"
const KEY_FILE = "06_server_key.pem"

const ROOT_DIR = "."
const BACKSLASH = "/"

const MESSAGE_TERMINATOR = '\n'
const MESSAGE_WHITESPACE = ' '

func main() {
	dir := GetDir(os.Args[1:]...)
	PrepareTls(CERT_FILE, KEY_FILE, func(c *tls.Config) {
		HandleTlsConnections(PROTOCOL, ADDRESS, c, func(c net.Conn) {
			CommandLoop(c, func(m string) {
				LoadJsonFile(dir, strings.ToLower(m), func(j string, e error) {
					log.Printf("%v: requested file %v\n", c.RemoteAddr(), m)
					var x string
					if e == nil {
						x = strings.ReplaceAll(j, string(MESSAGE_TERMINATOR), "")
					}
					SendMessage(c, x)
				})
			})
		})
	})
}

func PrepareTls(c, k string, f func(*tls.Config)) {
	cert, e := tls.LoadX509KeyPair(c, k)
	if e != nil {
		log.Fatal(e)
	}

	f(&tls.Config {
		Certificates: []tls.Certificate{ cert },
	})
}

func HandleTlsConnections(p, a string, c *tls.Config, f func(net.Conn)) {
	if l, e := tls.Listen(p, a, c); e == nil {
		for {
			if c, e := l.Accept(); e == nil {
				go func(c net.Conn) {
					defer c.Close()
					log.Printf("%v: connected\n", c.RemoteAddr())

					f(c)
				}(c)
			} else {
				log.Println(e)
			}
		}
	} else {
		log.Println(e)
	}
}

func CommandLoop(c net.Conn, f func(string)) {
	for {
		s, e := ReceiveMessage(c)
		switch e {
		case nil:
			for _, m := range MessageTokens(s) {
				f(m)
			}
		case io.EOF:
			log.Printf("%v: client connection dropped\n", c.RemoteAddr())
			return
		default:
			log.Printf("%v: %v\n", c.RemoteAddr(), e)
			return
		}
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
			log.Printf("%v: unable to send message [%T] %v\n", c.RemoteAddr(), v, v)
		}
	}
}

func MessageTokens(s string) []string {
	return strings.Split(s, string(MESSAGE_WHITESPACE))
}

func LoadJsonFile(dir, path string, f func(string, error)) {
	j, e := ioutil.ReadFile(AddSuffix(dir + path, ".json"))
	f(string(j), e)
}

func GetDir(s ...string) (d string) {
	if len(s) > 0 {
		d = RemoveDuplicates(s[0], BACKSLASH)
	} else {
		d = ROOT_DIR
	}
	return AddSuffix(d, BACKSLASH)
}

func AddSuffix(v, s string) string {
	if !strings.HasSuffix(v, s) {
		return v + s
	}
	return v
}

func RemoveDuplicates(s, sep string) string {
	return strings.Join(strings.Split(s, sep), sep)
}
