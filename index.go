package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"strings"
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
		case "/google/":
        	resp, err := http.Get("http://www.google.com")
        	if err != nil {
            	panic(err)
        	}

        	defer resp.Body.Close()
        	
        	body, err := ioutil.ReadAll(resp.Body)
    		if err != nil {
         		panic(err)
    		}
   		 	w.Header().Set("content-type", "text/html;charset=utf-8")
            w.Write([]byte(body))

          
    		//	fmt.Fprintf(w,string(body))    
    	case "/youtube/":
        	resp, err := http.Get("http://www.youtube.com")
        	if err != nil {
            	panic(err)
        	}

        	defer resp.Body.Close()
        	
        	body, err := ioutil.ReadAll(resp.Body)
    		if err != nil {
         		panic(err)
    		}
   		 	w.Header().Set("content-type", "text/html;charset=utf-8")
            w.Write([]byte(body))

        case "/search":
            resp, err := http.Get("http://www.google.com" + r.URL.String() )
            if err != nil {
            	panic(err)
        	}

        	defer resp.Body.Close()
        	
        	body, err := ioutil.ReadAll(resp.Body)
    		if err != nil {
         		panic(err)
    		}
   		// 	w.Header().Set("content-type", "text/html;charset=utf-8")
            w.Write([]byte(body))	

         case "/url":
            resp, err := http.Get("http://www.google.com" + r.URL.String() )
            if err != nil {
            	panic(err)
        	}

        	defer resp.Body.Close()
        	
        	body, err := ioutil.ReadAll(resp.Body)
    		if err != nil {
         		panic(err)
    		}
   	//	 	w.Header().Set("content-type", "text/html;charset=utf-8")
            w.Write([]byte(body))	

        default:  // 
         //	case "/images/":
    	 	var str string
    	 	str = r.URL.String()
    //	 	fmt.Fprintf(w, "PATH:"+r.URL.Path+"\r\n")
    	 	if strings.HasSuffix(str, "png"){
    	 		str = "http://www.google.com" + str
    	 		resp, err := http.Get(str)
    	 		if err != nil {
            	panic(err)
        		}

        		defer resp.Body.Close()
        	
        		body, err := ioutil.ReadAll(resp.Body)
    			if err != nil {
         			panic(err)
    			}
   		 		w.Header().Set("content-type", "image/png")
            	w.Write([]byte(body))

			}

			if strings.HasSuffix(str, "css"){
    	 		str = "http://www.youtube.com" + str
    	 		resp, err := http.Get(str)
    	 		if err != nil {
            	panic(err)
        		}

        		defer resp.Body.Close()
        	
        		body, err := ioutil.ReadAll(resp.Body)
    			if err != nil {
         			panic(err)
    			}
   		 		w.Header().Set("content-type", "text/css")
            	w.Write([]byte(body))

            }

            if strings.HasSuffix(str, "gif"){
    	 		str = "http://www.youtube.com" + str
    	 		resp, err := http.Get(str)
    	 		if err != nil {
            	panic(err)
        		}

        		defer resp.Body.Close()
        	
        		body, err := ioutil.ReadAll(resp.Body)
    			if err != nil {
         			panic(err)
    			}
   		 		w.Header().Set("content-type", "image/gif")
            	w.Write([]byte(body))

            }

 			if strings.HasSuffix(str, "jpg"){
    	 		str = "http://www.youtube.com" + str
    	 		resp, err := http.Get(str)
    	 		if err != nil {
            	panic(err)
        		}

        		defer resp.Body.Close()
        	
        		body, err := ioutil.ReadAll(resp.Body)
    			if err != nil {
         			panic(err)
    			}
   		 		w.Header().Set("content-type", "image/jpg")
            	w.Write([]byte(body))

            	}else {

				fmt.Fprintf(w, "METHOD:"+r.Method+"\r\n")
				fmt.Fprintf(w, "URL:\r\n")
				fmt.Fprintf(w, "PATH:"+r.URL.Path+"\r\n")
				fmt.Fprintf(w, "URL:"+r.URL.String()+"\r\n")

			}
        	

	}
}
