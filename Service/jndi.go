package Service

import (
	"Alphalog/Config"
	"bytes"
	"fmt"
	"log"
	"net"
	"strings"
)

type Result struct {
	Host   string
	Name   string
	Finger string
	Path   string
}

func JNDI() {

	log.Println("[JNDI]","Start fake reverse server")
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", Config.Init.PORT_JNDI))
	if err != nil {
		log.Fatalln("listen fail err: %s", err)
		return
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatalf("listen accept fail err: %s", err)
			continue
		}
		go acceptProcess(&conn)
	}
}

func acceptProcess(conn *net.Conn) {
	buf := make([]byte, 1024)
	num, err := (*conn).Read(buf)
	if err != nil {
		log.Fatalf("accept data reading err: %s", err)
		_ = (*conn).Close()
		return
	}
	hexStr := fmt.Sprintf("%x", buf[:num])
	// LDAP Protocol
	// https://ldap.com/ldapv3-wire-protocol-reference-bind
	if "300c020101600702010304008000" == hexStr {
		data := []byte{
			0x30, 0x0c, 0x02, 0x01, 0x01, 0x61, 0x07,
			0x0a, 0x01, 0x00, 0x04, 0x00, 0x04, 0x00,
		}
		_, _ = (*conn).Write(data)
		_, _ = (*conn).Read(buf)
		length := buf[8]
		pathBytes := bytes.Buffer{}
		for i := 1; i <= int(length); i++ {
			temp := []byte{buf[8+i]}
			pathBytes.Write(temp)
		}

		Host := (*conn).RemoteAddr().String()
		Path := pathBytes.String()
		jndi_log("ldap",Host,Path)
		_ = (*conn).Close()
		return
	}
	// RMI Protocol
	if checkRMI(buf) {
		data := []byte{
			0x4e, 0x00, 0x09, 0x31, 0x32,
			0x37, 0x2e, 0x30, 0x2e, 0x30,
			0x2e, 0x31, 0x00, 0x00, 0xc4, 0x12,
		}
		_, _ = (*conn).Write(data)
		_, _ = (*conn).Read(buf)
		_, _ = (*conn).Write([]byte{})
		_, _ = (*conn).Read(buf)
		var dataList []byte
		flag := false
		for i := len(buf) - 1; i >= 0; i-- {
			if buf[i] != 0x00 || flag {
				flag = true
				dataList = append(dataList, buf[i])
			}
		}
		var j int
		for i := 0; i < len(dataList); i++ {
			if int(dataList[i]) == i {
				j = i
			}
		}
		temp := dataList[0:j]
		pathBytes := &bytes.Buffer{}
		for i := len(temp) - 1; i >= 0; i-- {
			pathBytes.Write([]byte{dataList[i]})
		}

		Host := (*conn).RemoteAddr().String()
		Path := pathBytes.String()
		jndi_log("rmi",Host,Path)
		//ResultChan <- res

		_ = (*conn).Close()
		return
	}
	_ = (*conn).Close()
	return
}

// RMI Protocol Docs:
// https://docs.oracle.com/javase/9/docs/specs/rmi/protocol.html
func checkRMI(data []byte) bool {
	if data[0] == 0x4a &&
		data[1] == 0x52 &&
		data[2] == 0x4d &&
		data[3] == 0x49 {
		if data[4] != 0x00 {
			return false
		}
		if data[5] != 0x01 && data[5] != 0x02 {
			return false
		}
		if data[6] != 0x4b &&
			data[6] != 0x4c &&
			data[6] != 0x4d {
			return false
		}
		lastData := data[7:]
		for _, v := range lastData {
			if v != 0x00 {
				return false
			}
		}
		return true
	}
	return false
}

func jndi_log(log_type string,Host string,Path string){

	fuzz := strings.Split(Path,"/")
	log.Println(fuzz)
	subdomain_fuzz := fuzz[0] + "."+Config.Init.Domain+"."
	subdomain_flag,_,subdomain := RDB.Check_subdomain(subdomain_fuzz)
	if subdomain_flag{
		log.Println("["+log_type+"]",Host,Path)
		JNDI_log_data := RDB.Log_data(log_type,subdomain,Host,fuzz[1])
		RDB.PUSH(subdomain+"log",JNDI_log_data)
	}
}