package main

import (
	"log"
	"fmt"
	"net/http"
	"net/http/httputil"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!\r\n")
	fmt.Fprintf(w, "PATH:"+r.URL.Path+"\r\n")
	fmt.Fprintf(w, "SCHEME:"+r.URL.Scheme+"\r\n")
	fmt.Fprintf(w, "METHOD:"+r.Method+"\r\n")
	fmt.Fprintf(w, "HOST:"+r.Host+"\r\n")
	if r.URL.Path == "/dw"{
        remote, err := url.Parse("http://google.com")
        if err != nil {
                panic(err)
        }

        proxy := httputil.NewSingleHostReverseProxy(remote)
        fmt.Fprintf(w, "Proxying...")
        http.HandleFunc("/dw", handlerwww(proxy))
        if err != nil {
                panic(err)
        }
       

	}
}

func handlerwww(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
        return func(w http.ResponseWriter, r *http.Request) {
                log.Println(r.URL)
                w.Header().Set("X-Ben", "Rad")
                p.ServeHTTP(w, r)
        }
}