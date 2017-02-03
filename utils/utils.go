package utils

// transfer2go/utils - Go utilities for transfer2go
//
// Copyright (c) 2017 - Valentin Kuznetsov <vkuznet@gmail.com>

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"path/filepath"
	"runtime"
	"text/template"
)

// STATICDIR defines location of all static files
var STATICDIR string

// ListFiles function list files in given directory
func ListFiles(dir string) []string {
	var out []string
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println("Unable to read directory", dir, err)
		return nil
	}
	for _, f := range entries {
		if !f.IsDir() {
			out = append(out, f.Name())
		}
	}
	return out
}

// consume list of templates and release their full path counterparts
func fileNames(tdir string, filenames ...string) []string {
	flist := []string{}
	for _, fname := range filenames {
		flist = append(flist, filepath.Join(tdir, fname))
	}
	return flist
}

// ParseTmpl is a template parser with given data
func ParseTmpl(tdir, tmpl string, data interface{}) string {
	buf := new(bytes.Buffer)
	filenames := fileNames(tdir, tmpl)
	t := template.Must(template.ParseFiles(filenames...))
	err := t.Execute(buf, data)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

// Hash implements hash function for data, it returns a hash and number of bytes
func Hash(data []byte) (string, int64) {
	hasher := sha256.New()
	b, e := hasher.Write(data)
	if e != nil {
		log.Println("ERROR, Unable to write chunk of data via hasher.Write", e)
	}
	return hex.EncodeToString(hasher.Sum(nil)), int64(b)
}

// Stack helper function to return Stack
func Stack() string {
	trace := make([]byte, 2048)
	count := runtime.Stack(trace, false)
	return fmt.Sprintf("\nStack of %d bytes: %s\n", count, trace)
}

// ErrPropagate error helper function which can be used in defer ErrPropagate()
func ErrPropagate(api string) {
	if err := recover(); err != nil {
		log.Println("DAS ERROR", api, "error", err, Stack())
		panic(fmt.Sprintf("%s:%s", api, err))
	}
}

// ErrPropagate2Channel error helper function which can be used in goroutines as
// ch := make(chan interface{})
// go func() {
//    defer ErrPropagate2Channel(api, ch)
//    someFunction()
// }()
func ErrPropagate2Channel(api string, ch chan interface{}) {
	if err := recover(); err != nil {
		log.Println("DAS ERROR", api, "error", err, Stack())
		ch <- fmt.Sprintf("%s:%s", api, err)
	}
}

// GoDeferFunc helper function to run any given function in defered go routine
func GoDeferFunc(api string, f func()) {
	ch := make(chan interface{})
	go func() {
		defer ErrPropagate2Channel(api, ch)
		f()
		ch <- "ok" // send to channel that we can read it later in case of success of f()
	}()
	err := <-ch
	if err != nil && err != "ok" {
		panic(err)
	}
}

// FindInList helper function to find item in a list
func FindInList(a string, arr []string) bool {
	for _, e := range arr {
		if e == a {
			return true
		}
	}
	return false
}

// InList helper function to check item in a list
func InList(a string, list []string) bool {
	check := 0
	for _, b := range list {
		if b == a {
			check += 1
		}
	}
	if check != 0 {
		return true
	}
	return false
}

// MapKeys helper function to return keys from a map
func MapKeys(rec map[string]interface{}) []string {
	keys := make([]string, 0, len(rec))
	for k := range rec {
		keys = append(keys, k)
	}
	return keys
}

// List2Set helper function to convert input list into set
func List2Set(arr []string) []string {
	var out []string
	for _, key := range arr {
		if !InList(key, out) {
			out = append(out, key)
		}
	}
	return out
}

// HostIP provides a list of host IPs
func HostIP() []string {
	var out []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Println("ERROR unable to resolve net.InterfaceAddrs", err)
	}
	for _, addr := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				out = append(out, ipnet.IP.String())
			}
			if ipnet.IP.To16() != nil {
				out = append(out, ipnet.IP.String())
			}
		}
	}
	return List2Set(out)
}