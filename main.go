package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type flags struct {
	url     url.URL
	path    string
	content []byte
}

var Dat *flags

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU()*2 + 1)

	f, err := getFlags()
	if err != nil {
		log.Fatalf("flags parsing fail: %v", err)
	}

	Dat = &f

	fs := http.FileServer(http.Dir(f.path))
	http.Handle("/file/", http.StripPrefix("/file/", fs))

	http.HandleFunc("/", handler)

	err = http.ListenAndServeTLS(getPort(f.url), "cert.pem", "key.pem", nil)
	if err != nil {
		log.Fatalf("ListenAndServe: ", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {

	if Dat.content == nil {
		fmt.Fprint(w, r.Host+","+strconv.FormatInt(time.Now().UnixNano(), 10))
	} else {
		fmt.Fprint(w, string(Dat.content))
	}
}

func getFlags() (flags, error) {

	u := flag.String("url", "http://localhost:10443", "server")

	p := flag.String("path", "/home/ubuntu/file_server", "file folder")

	fn := flag.String("file", "content.html", "file name")

	flag.Parse()

	ur, err := url.Parse(*u)
	if err != nil {
		log.Printf("url parse err: %v", err)
		return flags{}, err
	}

	body, err := ioutil.ReadFile(*fn)
	if err != nil {
		log.Printf("ioutil.ReadFile err: %v", err)
		body = nil
	}

	return flags{*ur, *p, body}, nil
}

func getPort(u url.URL) string {

	r := u.Host

	if n := strings.Index(r, ":"); n != -1 {
		r = r[n:]
	} else {
		r = ":10443"
	}

	return r
}
