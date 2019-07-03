package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"io"
	"strings"
	"strconv"
	"bytes"
	"compress/gzip"
	"database/sql"
	"bufio"
    _ "github.com/go-sql-driver/mysql"
)
const zhost string = `https://v2ray.14065567.now.sh/`

func Handler(w http.ResponseWriter, r *http.Request) {
	var (
		url =``
		realhost =``
	)

	switch r.URL.Path{
		case `/manager/`:
			db, err := sql.Open("mysql","zhujq:Juju1234@tcp(35.230.121.24:3316)/zeit")
			if err != nil {
				fmt.Fprintf(w, err.Error() )
				return
			}
			defer db.Close()
			err = db.Ping()
			if err != nil {
				fmt.Fprintf(w, err.Error() )
				return
			}

			rows, err := db.Query("select * from visits where to_days(visitime) = to_days(now()) order by id desc;")
			if err != nil {
				fmt.Println(err.Error() )	
				return
			}
			defer rows.Close()
			
			fmt.Fprintf(w,`%s`,`<html><head></head><body><table width="100%" border="1" align="center"  font-size:9px color="#660066" style="word-break:break-all; word-wrap:break-all;"><tr><th width="2%">id</th><th width="10%">datetime</th><th idth="3%">Method</th><th width="20%">URL</th><th width="30%">Req_Head</th><th width="5%">RSP_STATUS</th><th width="25%">RSP_HEAD</th><th width="5%%">RSP_LENGTH</th></tr>`)
			var (
				visitid = 0
				visitime =``
				visitmethod =``
				visiturl =``
				visithead =``
				rsp_status =``
				rsp_head =``
				rsp_length = 0
			)
			for rows.Next() {	
				if err = rows.Scan(&visitid,&visitime,&visitmethod,&visiturl,&visithead,&rsp_status,&rsp_head,&rsp_length); err != nil {
					fmt.Println(err.Error() )	
				}
				fmt.Fprintf(w,`<tr><td width="2%%">%d</td><td width="10%%">%s</td><td width="3%%">%s</td><td width="20%%">%s</td><td width="30%%">%s</td><td width="5%%">%s</td><td width="25%%">%s</td><td width="5%%">%d</td></tr>`,visitid,visitime,visitmethod,visiturl,visithead,rsp_status,rsp_head,rsp_length)
			}
				
			if err = rows.Err(); err != nil {
				fmt.Println(err.Error() )	
			}
			fmt.Fprintf(w,"%s","</table></body>")
			return

		case `/google/`:    //google入口
			url = zhost + `www.google.com/`
			realhost = `www.google.com`
			http.Redirect(w, r, url, 307)
			return
			  
    	case `/youtube/`:   //youtube入口
			url = zhost + `www.youtube.com/`
			realhost = `www.youtube.com`
			http.Redirect(w, r, url, 307)
			return	
		
		case `/watch`:   //youtube入口
			url = `https://www.youtube.com`+ r.URL.String() 
			realhost = `www.youtube.com`

        default:    //  经google、youtube入口后重新返回的网址和直接带上的域名的处理，分离出真实主机名称 
    	 	var str string
			str = r.URL.String()
			str = strings.TrimLeft(str,`/`)
			realhost = string([]byte(str)[0:strings.Index(str,`/`)])  //去掉首位的/后截取host
			if realhost == `signin` {             //google的signin目录暂时无法带上accounts.google.com
				realhost = `accounts.google.com`
				str = realhost + `/` + str
			}

			if realhost == `youtubei` || realhost == `yts` || realhost == `results` {      //youtube的youtubei yts results目录暂时无法带上www.youtube.com
				realhost = `www.youtube.com`
				str = realhost + `/` + str
			}
		
    	 	if realhost == ``{
				fmt.Fprintf(w, `Failed to handle RequestUrl:`+str+`\r\n`)
				return
			}
			if r.URL.Scheme == ``{
				r.URL.Scheme = `https`
			}
			url = r.URL.Scheme + `://` + str
        	
	}

	if toredirect(realhost){             //判断如果是国内域名，则指示重定向
	//	fmt.Println(r.Method,` URL:`+url,` LocalRealHost:`,realhost)	
		http.Redirect(w, r, url, 307)
		return
	}
	
	client := &http.Client{}
	req, err := http.NewRequest(r.Method, url, nil)
	if err != nil {
        panic(err)
    }
	req.Header = r.Header     //删除请求头压缩选项，否则无法对返回的文本的链接内容进行处理,20190625 调用compress/gzip进行压缩和解压缩,且只用gzip
	if  strings.Contains(string(req.Header.Get(`Accept-Encoding`)),`gzip`){
		req.Header.Set(`Accept-Encoding`,`gzip`)  
	}else {
		req.Header.Del(`Accept-Encoding`)   
	}
	strreferer := string(req.Header.Get("Referer"))     //修改referer
	if strings.Contains(strreferer,`v2ray.14065567.now.sh/`){
		strreferer = strings.Replace(strreferer,`v2ray.14065567.now.sh/`,``,-1)
		req.Header.Set(`Referer`,strreferer)
	}
	req.Header.Set(`Host`,realhost)    // 设置请求头的真实host
	//删除zeit添加的头域
	req.Header.Del(`X-Forwarded-For`)
	req.Header.Del(`X-Zeit-Co-Forwarded-For`)
	req.Header.Del(`X-Real-Ip`)
	req.Header.Del(`X-Now-Trace`)
	req.Header.Del(`X-Forwarded-Host`)
	req.Header.Del(`X-Now-Deployment-Url`)
	req.Header.Del(`X-Forwarded-Proto`)
	req.Header.Del(`X-Now-Id`)

	
	req.Body = r.Body   //加入POST时的Body
	req.Form = r.Form
	req.PostForm = r.PostForm
	req.MultipartForm = r.MultipartForm

	resp, err := client.Do(req)
	
	fmt.Println(r.Method,` URL:`+url,`resp len,status,type,Enc:`,strconv.FormatInt(resp.ContentLength,10),resp.Status,resp.Header.Get(`content-type`),resp.Header.Get(`Content-Encoding`))	//记录访问记录
	s, _ := ioutil.ReadAll(r.Body)
	fmt.Println(`Request Body:`)
	fmt.Println(s)
    if err != nil {
        panic(err)
    }

	db, err := sql.Open("mysql","zhujq:Juju1234@tcp(35.230.121.24:3316)/zeit")
	if err == nil {
		err = db.Ping()
		if err == nil {
			reqhead := ``
			for k, _ := range r.Header {
				reqhead += k
				reqhead += `:`
				reqhead += r.Header.Get(k)
				reqhead += `<br>`
			}
			reqhead =  strings.Replace(reqhead,`"`,`\"`,-1) //存入mysql时要把“转义
			rsphead := ``
			for k, _ := range resp.Header {
				rsphead += k
				rsphead += `:`
				rsphead += resp.Header.Get(k)
				rsphead += `<br>`
			}
			rsphead =  strings.Replace(rsphead,`"`,`\"`,-1)
			var insertsql = `insert into visits(method,url,head,rsp_status,rsp_head,rsp_legnth) values(`+`"` + r.Method +`","` + url +`","`+ reqhead +`","` + resp.Status +`","` + rsphead + `","` + strconv.FormatInt(resp.ContentLength,10)+`");`
		//	fmt.Println(insertsql)	
			_,err := db.Exec(insertsql)
			if err != nil{
				fmt.Println(err.Error() )	
			}

		}else{
			fmt.Println(err.Error() )	
		}
		
	}else{
		fmt.Println(err.Error() )	
	}
	defer db.Close()

    defer resp.Body.Close()
        	
    for k, _ := range resp.Header{
   		w.Header().Set(k,resp.Header.Get(k))
	}
	
	if strings.Contains(string(resp.Header.Get(`content-type`)),`text`) {  //返回文本类型时的处理

		body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
        	panic(err)
		}
	
		if resp.StatusCode == 200 && strings.Contains(string(resp.Header.Get(`content-type`)),`text`){   //只有当返回200和文本类型时进行链接处理
			if len(body) == 0 {
				fmt.Println(`resp is empty`)
				return
			}
		
			if resp.Header.Get(`Content-Encoding`) == `gzip`{
				body,err = gzipdecode(body)
				if err != nil {
					panic(err)
				}
			}
		
			modifiedrsp := []byte{}
			tomodifystr := ``
			for _,v := range body {
				if string(v) == `<` {
					if len(tomodifystr) >0 {
						tomodifystr = modifylink(tomodifystr,realhost)
						for _,vv := range tomodifystr {
							modifiedrsp = append(modifiedrsp,byte(vv))
						}
						tomodifystr = ``
					}
					tomodifystr += string(v)
				}else if string(v) == `>` {
					tomodifystr += string(v)
					tomodifystr = modifylink(tomodifystr,realhost)
					for _,vv := range tomodifystr {
						modifiedrsp = append(modifiedrsp,byte(vv))
					}
					tomodifystr = ``
				}else{
					tomodifystr += string(v)
				}
			}

			body = modifiedrsp
		
		
			if resp.Header.Get("Content-Encoding") == "gzip" {    //如果resp指示压缩，还需要对解开的处理后的内容重新压缩
				body,err = gzipencode(body)
				if err != nil {
					panic(err)
				}

			}

		}

		w.Write([]byte(body)) 
	}else {   //返回非文本类型，用stream模式处理
		reader := bufio.NewReader(resp.Body)
		caches := make([]byte, 10240) 
		for {
			_, err :=reader.Read(caches)
			if err != nil {
				if err == io.EOF {
					break
				}else{
					panic(err)
				}
			} 
			w.Write(caches)
		}
	}       
}

func modifylink(s string,realhost string) string{
	if s == ""{
		return s
	}

	tempstr := s	                                                     
	tempstr =  strings.Replace(tempstr,`href="https://`, `href="` + zhost, -1) //先改https，否则会重复改	                                 
	tempstr =  strings.Replace(tempstr,`href= "https://`,`href="` + zhost,-1) //youtube上发现有href=空格情况
	tempstr =  strings.Replace(tempstr,`href="http://`,`href="` + zhost,-1)
	tempstr =  strings.Replace(tempstr,`href="//`,`href="` +zhost,-1)  // href="//  后是绝对路径
	tempstr =  strings.Replace(tempstr,`href="/`,`href="` + zhost + realhost + "/",-1)   // href="/ 后是相对路径
		
	tempstr =  strings.Replace(tempstr,`<a href="https://`,`<a href="` + zhost,-1)
	tempstr =  strings.Replace(tempstr,`<a href="http://`,`<a href="` + zhost,-1)
	tempstr =  strings.Replace(tempstr,`<a href="//`,`<a href="` + zhost,-1)
	tempstr =  strings.Replace(tempstr,`<a href="/`,`<a href="` + zhost + realhost + "/",-1)
	
	tempstr =  strings.Replace(tempstr,`action="https://`,`action="` + zhost,-1)
	tempstr =  strings.Replace(tempstr,`action="http://`,`action="` + zhost,-1)
	tempstr =  strings.Replace(tempstr,`action="//`,`action="` + zhost,-1)
	tempstr =  strings.Replace(tempstr,`action="/`,`action="` + zhost + realhost + "/",-1)

	tempstr =  strings.Replace(tempstr,`location='/`,`location='` + zhost + realhost + "/",-1)

	tempstr =  strings.Replace(tempstr,`src="https://`,`src="` + zhost,-1)
	tempstr =  strings.Replace(tempstr,`src="http://`,`src="` + zhost,-1)
	tempstr =  strings.Replace(tempstr,`src="//`,`src="` + zhost,-1)
	tempstr =  strings.Replace(tempstr,`src="/`,`src="` +zhost  + realhost+ "/",-1)
		
	tempstr =  strings.Replace(tempstr,`srcset="https://`,`srcset="` + zhost,-1)
	tempstr =  strings.Replace(tempstr,`srcset="http://`,`srcset="` + zhost,-1)
	tempstr =  strings.Replace(tempstr,`srcset="//`,`srcset="` + zhost,-1)
	tempstr =  strings.Replace(tempstr,`srcset="/`,`srcset="` + zhost  + realhost+ "/",-1)
	
	tempstr =  strings.Replace(tempstr,`<meta content="https://`,`<meta content="` + zhost,-1)
	tempstr =  strings.Replace(tempstr,`<meta content="http://`,`<meta content="` + zhost,-1)
	tempstr =  strings.Replace(tempstr,`<meta content="//`,`<meta content="` + zhost,-1)
	tempstr =  strings.Replace(tempstr,`<meta content="/`,`<meta content="` + zhost  + realhost + "/",-1)

	tempstr =  strings.Replace(tempstr,`<iframe src="https://`,`<iframe src="` + zhost,-1)
	tempstr =  strings.Replace(tempstr,`<iframe src="http://`,`<iframe src="` + zhost,-1)
	
	tempstr =  strings.Replace(tempstr,`itemtype="https://`,`itemtype="` + zhost,-1)
	tempstr =  strings.Replace(tempstr,`itemtype="http://`,`itemtype="` + zhost,-1)

	//script的url修正
	tempstr = strings.Replace(tempstr,`url(https://`,`url(` + zhost, -1)
	tempstr = strings.Replace(tempstr,`url('https://`,`url('` +zhost, -1)
	tempstr = strings.Replace(tempstr,`url(//`,`url(` +zhost + realhost + `/`, -1)
	tempstr = strings.Replace(tempstr,`url("//`,`url("` +zhost, -1)
	tempstr = strings.Replace(tempstr,`url(/`,`url(` +zhost + realhost + `/`, -1)
	tempstr = strings.Replace(tempstr,`url=https://`,`url=` +zhost,-1)
	tempstr = strings.Replace(tempstr,`s='/images`,`s='` +zhost + realhost + `/images`, -1)
	tempstr = strings.Replace(tempstr,`http:\/\/`,`https:\/\/` +`v2ray.14065567.now.sh` + `\/`,-1)
	tempstr = strings.Replace(tempstr,`https:\/\/`,`https:\/\/` +`v2ray.14065567.now.sh` + `\/`,-1)
	tempstr = strings.Replace(tempstr,`"url":"https://`,`"url":"` +zhost ,-1)
	tempstr = strings.Replace(tempstr,`"url":"/`,`"url":"` +zhost + realhost + `/`, -1)
	tempstr = strings.Replace(tempstr,`='/`,`='` +zhost + realhost + `/`, -1)
	

	tempstr = strings.Replace(tempstr,`v2ray.14065567.now.sh/v2ray.14065567.now.sh`,`v2ray.14065567.now.sh`, -1)	 //有可能重复修改	
	tempstr = strings.Replace(tempstr,`v2ray.14065567.now.sh/https://v2ray.14065567.now.sh`,`v2ray.14065567.now.sh`, -1)
	
	
	return tempstr

}

func toredirect(s string) bool{
	if strings.HasSuffix(s, ".cn"){
		return true
	}
	localurls := []string{"baidu","taobao","sina","163.com","tmall","jd.com","sohu","qq.com","ifeng.com","qunae.com","toutiao.com","alipay.com","ctrip.com","weibo.com","zhihu","csdn.net","58.com","youku","iyiyi.com","bilibili"}
	for _, localurl := range localurls {
		if strings.Contains(s,localurl){
			return true
		}
	}
	
	return false
}

func gzipencode(in []byte) ([]byte, error) {
    var (
        buffer bytes.Buffer
        out    []byte
		err    error
       	)
        writer := gzip.NewWriter(&buffer)
        _, err = writer.Write(in)
        if err != nil {
            writer.Close()
            return out, err
        }
        err = writer.Close()
        if err != nil {
            return out, err
        }
   	     return buffer.Bytes(), nil
}
	
func gzipdecode(in []byte) ([]byte, error) {
        reader, err := gzip.NewReader(bytes.NewReader(in))
        if err != nil {
            var out []byte
            return out, err
        }
        defer reader.Close()
        return ioutil.ReadAll(reader)
}


