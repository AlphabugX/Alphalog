package Data

import (
	"Alphalog/Config"
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"strings"
	"time"
)

type RedisDB struct {
	Client *redis.Client
}

var ctx = context.Background()

func (RDB *RedisDB) GET(key string) (string, error) {
	value, err := RDB.Client.Get(ctx, key).Result()
	return value, err
}
func (RDB *RedisDB) SET(key string, value string) (string, error) {
	result, err := RDB.Client.Set(ctx, key, value, 0).Result()
	RDB.Client.Expire(ctx, key, 24*time.Hour)
	return result, err
}
func (RDB *RedisDB) Exist(key string) (int64, error) {
	key = strings.Replace(key, "\n", "", -1)
	log.Println("[Redis][Exist][key]", key)
	result, err := RDB.Client.Exists(ctx, key).Result()
	if err != nil {
		log.Println(err)
	}
	return result, err
}
func (RDB *RedisDB) PUSH(key string, value string) error {
	_, err := RDB.Client.RPush(ctx, key, value).Result()
	RDB.Client.Expire(ctx, key, 24*time.Hour)
	return err
}
func (RDB *RedisDB) RANGE(key string) (interface{}, error) {
	result, err := RDB.Client.LRange(ctx, key, 1, -1).Result()
	return result, err
}
func (RDB *RedisDB) Getkey(key string) (string, error) {
	result, err := RDB.Client.LIndex(ctx, key, 0).Result()
	return result, err
}
func RedisInit() RedisDB {
	RClient := RedisDB{Client: Config.RClient}
	log.Println("[Redis]", "Start redis database init")
	return RClient
}
func (RDB *RedisDB) Log_data(types string, subdomain string, ip string, reqbody interface{}) (result string) {
	timeNow := time.Unix(time.Now().Unix(), 0)
	logs := map[string]interface{}{
		"type":      types,
		"subdomain": subdomain,
		"ip":        ip,
		"reqbody":   reqbody,
		"time":      timeNow.String(),
	}
	logs2json, _ := json.Marshal(logs)
	result = string(logs2json)
	return result
}

func (RDB *RedisDB) Check_subdomain(subdomain string) (bool bool, check_subdomain string, fuzz_domain string) {
	subdomain_len := len(Config.Init.Domain) + 6
	if len(subdomain) < subdomain_len {
		return false, check_subdomain, fuzz_domain
	}
	fuzz_domain = subdomain[len(subdomain)-subdomain_len:] // xxxx.fuzz.red.

	check_subdomain = "" //	=> Alphabug.xxxx.fuzz.red.

	flag := false
	if len(subdomain) == len(fuzz_domain) {
		flag = true
		check_subdomain = subdomain
	} else if len(subdomain) > len(fuzz_domain) {
		// => xxxx.fuzz.red
		// => Alphabug.xxxx.fuzz.red
		check_subdomain = subdomain[len(subdomain)-len(fuzz_domain):]
		check_dot := subdomain[0 : len(subdomain)-subdomain_len] // => Alphabug.
		check_dot = check_dot[len(check_dot)-1:]                 // => .
		if check_dot == "." {
			flag = true
		}
	}
	ok, _ := RDB.Exist(fuzz_domain + "log")
	if ok == 1 && flag {
		//log.Println("check_subdomain",check_subdomain)
		log.Println("[Check_subdomain]", subdomain, fuzz_domain, "True")
		return true, check_subdomain, fuzz_domain
	}
	log.Println("[Check_subdomain]", subdomain, fuzz_domain, "False")
	return false, check_subdomain, fuzz_domain
}
