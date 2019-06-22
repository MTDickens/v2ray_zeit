package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path{
		case "/":
			fmt.Fprintf(w, "Welcome you,your info:\r\n")
			fmt.Fprintf(w, "METHOD:"+r.Method+"\r\n")
			fmt.Fprintf(w, "URL:\r\n")
			fmt.Fprintf(w, "PATH:"+r.URL.Path+"\r\n")
			fmt.Fprintf(w, "SCHEME:"+r.URL.Scheme+"\r\n")
			fmt.Fprintf(w, "HOST:"+r.URL.Host+"URL-End\r\n")
    		fmt.Fprintf(w, "Proto:"+r.Proto+"\r\n")
			fmt.Fprintf(w, "HOST:"+r.Host+"\r\n")
			fmt.Fprintf(w, "RequestUrl:"+r.RequestURI+"\r\n")
		case "/dw":
        	resp, err := http.Get("http://www.google.com")
        	if err != nil {
            	panic(err)
        	}

        	defer resp.Body.Close()
        	body, err := ioutil.ReadAll(resp.Body)
    		if err != nil {
          		panic(err)
    		}
 
    		fmt.Fprintf(w,string(body))    
    	case "/youtube"
    	    w, err := http.Get("http://www.youtube.com")
    	    if err != nil {
            	panic(err)
        	}

        	defer w.Body.Close()
        	if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}

        default:

	}
}
