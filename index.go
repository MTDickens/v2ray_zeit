package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	var url, realhost string
	if r.URL.Scheme == ""{
		r.URL.Scheme = "http"
	}
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
			return

		case "/google/":    //google入口
			url = "http://www.google.com"
			realhost = "www.google.com"
			  
    	case "/youtube/":   //youtube入口
			url = "http://www.youtube.com"
			realhost = "www.youtube.com"
		
		case "/search":     //google search入口，由于暂时无法带上真实主机名导致
            url = "http://www.google.com" + r.URL.String() 
            realhost = "www.google.com"

		// case "/url":
        // 		url = "http://www.google.com" + r.URL.String() 
			    	

        default:    //  经google、youtube入口后重新返回的网址的处理，分离出真实主机名称 
    	 	var str string
    	 	str = r.URL.String()
			realhost = string([]byte(str)[1:strings.Index(str,"/")])
			url = str
    	 	if realhost == ""{
				fmt.Fprintf(w, "Failed to handle RequestUrl:"+str+"\r\n")
			}
        	
	}

	client := &http.Client{}
	req, err := http.NewRequest(r.Method, url, nil)
	req.Header = r.Header
	if err != nil {
        panic(err)
    }

    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }

    defer resp.Body.Close()
        	
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }
   		 //	w.Header().Set("content-type", "text/html;charset=utf-8")
   	for k, _ := range resp.Header{
   		w.Header().Set(k,resp.Header.Get(k))
   	}
	
//	if strings.Contains(string(resp.Header.Get("content-type")),"text/html"){  //考虑对所有返回的文本进行链接修正，指回zeit

		olds := "<a href=\"/"
		news := "<a href=" + "\""+ "https://v2ray.14065567.now.sh/" + realhost + "/"
		body = []byte(strings.ReplaceAll(string(body),olds,news))

		olds = "src=\"/"
		news = "src=" + "\""+ "https://v2ray.14065567.now.sh/" + realhost+ "/"
		body = []byte(strings.ReplaceAll(string(body),olds,news))

		olds = "href=\"http://"
		news = "href=" + "\""+ "https://v2ray.14065567.now.sh/" 
		body = []byte(strings.ReplaceAll(string(body),olds,news))

		olds = "href=\"https://"
		news = "href=" + "\""+ "https://v2ray.14065567.now.sh/" 
		body = []byte(strings.ReplaceAll(string(body),olds,news))

		olds = "<meta content=\"https://"
		news = "<meta content=" + "\""+ "https://v2ray.14065567.now.sh/" 
		body = []byte(strings.ReplaceAll(string(body),olds,news))

		olds = "<meta content=\"/"
		news = "<meta content=" + "\""+ "https://v2ray.14065567.now.sh/" + realhost + "/"
		body = []byte(strings.ReplaceAll(string(body),olds,news))
		

//	}
	fmt.Println(r.Method," URL:"+url," RealHost:",realhost,resp.Header.Get("content-type"))		
			
    w.Write([]byte(body))

          
    //	fmt.Fprintf(w,string(body))  
}
