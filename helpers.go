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
)

const ADDRESS = ":1024"
const HTTPS_URL = "https://localhost:1024/"

const ROOT_DIR = "."
const BACKSLASH = "/"

const HTML = ".html"
const JSON = ".json"

const MESSAGE_TERMINATOR = '\n'
const MESSAGE_WHITESPACE = ' '

type Person struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  string
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

func ReceiveMessage(c net.Conn) (b []byte, e error) {
	if b, e = bufio.NewReader(c).ReadBytes(MESSAGE_TERMINATOR); e == nil {
		b = b[:len(b)-1]
	}
	return
}

func Tokens(b []byte) []string {
	return strings.Split(string(b), string(MESSAGE_WHITESPACE))
}

func ServerUrl(s ...string) string {
	url := []string{HTTPS_URL}
	return GetPath(append(url, s...)...)
}

func LoadFile(t, p string, f func(string, error)) {
	j, e := ioutil.ReadFile(AddSuffix(p, t))
	f(string(j), e)
}

func GetPath(n ...string) string {
	for i, v := range n {
		n[i] = strings.Trim(v, BACKSLASH)
	}
	return strings.Join(n, BACKSLASH)
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

func ForEachRecord[T any](b []byte, f func(T)) {
	r := []T{}
	LogErrors(json.Unmarshal(b, &r), func() {
		for _, v := range r {
			f(v)
		}
	})
}

func Parallelize(s []string, f func(string)) {
	var w sync.WaitGroup

	for _, n := range s {
		w.Add(1)
		go func(n string) {
			defer w.Done()

			f(n)
		}(n)
	}
	w.Wait()
}

func LogErrors(e error, f ...func()) {
	fs, fe := BodyAndTail(f)
	if e == nil {
		for _, x := range fs {
			x()
		}
	} else {
		log.Println(e)
		if fe != nil {
			fe()
		}
	}
}

func BodyAndTail[T any](s []T) ([]T, T) {
	var rt T
	l := len(s) - 1
	b, t := s[:l], s[l:]
	if len(b) < len(t) {
		b = append(b, t[0])
		t = t[1:]
	}
	if len(t) > 0 {
		rt = t[0]
	}
	return b, rt
}
