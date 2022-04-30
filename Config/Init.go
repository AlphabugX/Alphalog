package Config

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type Config struct {
	Database interface{}
	Domain string
	IP_DNS string
	IP_JNDI string
	PORT_JNDI string
	PORT_HTTP string

}
var Init Config
var RClient *redis.Client
func Initialization() {
	Init.InitConfig()
}
func (Init *Config)InitConfig()  {
	conf := yamlread()
	Init.Domain = conf["domain"]
	Init.IP_DNS = conf["IP_DNS"]
	Init.IP_JNDI = conf["IP_JNDI"]
	Init.PORT_JNDI = conf["PORT_JNDI"]
	Init.PORT_HTTP = conf["PORT_HTTP"]
	database := strings.Split(conf["database"],":")
	switch database[0] {
	case "redis":
		Init.Database = map[interface{}]string{
			"type":database[0],
			"ip": database[1],
			"port":database[2],
			"password":database[3],
		}
		if Init.redis_config() {
			log.Println("[Config] Redis Initialization succeeded")
		}else{
			log.Fatalln("[Config] Redis Initialization Failed")
		}
	default:
		Init.Database = nil
	}
}


func (Init Config)redis_config()bool {
	RedisYaml := Init.Database.(map[interface{}]string)
	RDB := redis.NewClient(&redis.Options{
		Addr:		RedisYaml["ip"]+":"+RedisYaml["port"],
		Password:	RedisYaml["password"],
		DB:			0,
	})
	var ctx,cancel = context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()
	_, err := RDB.Ping(ctx).Result()
	if err == nil{
		RClient = RDB
		return true
	}else {
		return false
	}
}
func yamlread()(conf map[interface{}]string)  {
	yamlfile,err:=os.Open("config.yaml")
	if err != nil {
		log.Fatalln("[Config] config.yaml file does not exist.")
	}
	yamldata,_:= ioutil.ReadAll(yamlfile)
	conf = make(map[interface{}]string)
	b := yaml.Unmarshal(yamldata,&conf)
	if b != nil {
		log.Fatalln("[Config] config.yaml initialization failed.")
		return nil
	}
	return conf
}