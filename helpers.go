package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"
	"sync"
	"syscall"
)

type Person struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  string
}

func GetSocketOption(c *net.UDPConn, p int) (r int) {
	fd, _ := c.File()
	r, _ = syscall.GetsockoptInt(int(fd.Fd()), syscall.SOL_SOCKET, p)
	return
}

func SendMessage(c net.Conn, s ...any) {
	for _, v := range s {
		switch v := v.(type) {
		case []byte:
			c.Write(v)
			c.Write([]byte(string('\n')))
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

func ReceiveMessage(c net.Conn) (b []byte, e error) {
	if b, e = bufio.NewReader(c).ReadBytes('\n'); e == nil {
		b = b[:len(b)-1]
	} else {
		log.Println(e)
	}
	return
}

func ReadStream(c net.Conn) (r []byte) {
	var e error
	var n int
	m := make([]byte, 1024)
	if n, e = c.Read(m); e == nil {
		r = m[:n]
	} else {
		log.Println(e)
	}
	return
}

func Tokens(b []byte) []string {
	return strings.Split(string(b), string(' '))
}

func ServerUrl(s ...string) string {
	return GetPath(append([]string{"https://localhost:1024/"}, s...)...)
}

func LoadFile(t, p string) (b []byte) {
	var e error
	if b, e = ioutil.ReadFile(AddSuffix(p, t)); e != nil {
		log.Println(e)
	}
	return b
}

func GetPath(n ...string) string {
	for i, v := range n {
		n[i] = strings.ToLower(strings.Trim(v, "/"))
	}
	return strings.Join(n, "/")
}

func GetDir(s ...string) (d string) {
	if len(s) > 0 {
		d = RemoveDuplicates(s[0], "/")
	} else {
		d = "."
	}
	return AddSuffix(d, "/")
}

func DeleteAll[T comparable](s []T, p T) (r []T) {
	for _, v := range s {
		if v != p {
			r = append(r, v)
		}
	}
	return
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

func ForEachRecord[T any](b []byte, f func(T)) {
	r := []T{}
	if e := json.Unmarshal(b, &r); e == nil {
		for _, v := range r {
			f(v)
		}
	} else {
		log.Println(e)
	}
}

func Parallelize[T any](s []T, f func(T)) {
	var w sync.WaitGroup

	for _, n := range s {
		w.Add(1)
		go func(n T) {
			defer w.Done()

			f(n)
		}(n)
	}
	w.Wait()
}
