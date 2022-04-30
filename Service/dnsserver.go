package Service


import (
	"Alphalog/Config"
	"Alphalog/Data"
	"fmt"
	"github.com/miekg/dns"
	"log"
	"strings"
)

var RDB = Data.RedisDB{}

func Dnsserver(domain string) {
	log.Println("[DNS]","Domain:",domain)
	dns.HandleFunc(domain+".", handleDnsRequest)
	// 解析的一级域名 例如：alphabug.cn 或者 dnslog.cn
	RDB = Data.RedisInit()
	server := &dns.Server{Addr: ":53", Net: "udp"}
	log.Println("[DNS]","Initialization succeeded")
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}

func parseQuery(RemoteAddr string,m *dns.Msg) {


	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:

			q.Name = strings.ToLower(q.Name)
			//timeNow := time.Unix(time.Now().Unix(),0)
			log.Println("[DNSQuery]",q.Name,RemoteAddr)

			ip :=""
			check_subdomain_flag,_,subdomain := RDB.Check_subdomain(q.Name)
			if q.Name == "www."+Config.Init.Domain+"." || q.Name == Config.Init.Domain+"."{
				ip = Config.Init.IP_DNS
			}else if check_subdomain_flag || q.Name == "jndi."+Config.Init.Domain+"." {
				ip = Config.Init.IP_JNDI
			}else{
				ip = "127.0.0.1"
			}

			if  check_subdomain_flag{
				//log.Println("[Check]",check_subdomain_strings)
				//dnslog :=q.Name +"|"+ strings.Split(RemoteAddr,":")[0]+"|" +timeNow.String()
				//A :=
				dnslog := RDB.Log_data("dns",q.Name,strings.Split(RemoteAddr,":")[0],"")
				RDB.PUSH(subdomain+"log",dnslog)
			}

			rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
			if err == nil {
				m.Answer = append(m.Answer, rr)
			}

		}
	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)

	m.SetReply(r)
	m.Compress = false
	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(w.RemoteAddr().String(),m)
	}
	w.WriteMsg(m)

}

