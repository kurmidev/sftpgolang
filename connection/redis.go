package connection

import (
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

var REDIS_HOST string = GoDotEnvVariable("REDIS_HOST")
var REDIS_USERNAME string = GoDotEnvVariable("REDIS_USERNAME")
var REDIS_PASSWORD string = GoDotEnvVariable("REDIS_PASSWORD")
var REDIS_PORT string = GoDotEnvVariable("REDIS_PORT")
var REDIS_DB string = GoDotEnvVariable("REDIS_DB")

func RedisConnect() {

	//url := "redis://:@ll-prod-redis-server.fcgww7.ng.0001.aps1.cache.amazonaws.com:6379/10?protocol=3"
	url := fmt.Sprintf("redis://%s:%s@%s:%s/%s?protocol=3", REDIS_USERNAME, REDIS_PASSWORD, REDIS_HOST, REDIS_PORT, REDIS_DB)
	redisdb, err := redis.ParseURL(url)

	if err != nil {
		log.Fatal(err)
	}

	rdb = redis.NewClient(redisdb)
}

func GetRedisDb() *redis.Client {
	return rdb
}
