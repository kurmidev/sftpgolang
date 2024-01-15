package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/liquiloans/sftp/config"
	"github.com/liquiloans/sftp/connection"
	"github.com/liquiloans/sftp/utils"
	"github.com/redis/go-redis/v9"
)

var redisdb *redis.Client
var ctx = context.Background()

type RedisData struct {
	InvestorId int         `json:"investorid"`
	IfaId      int         `json:"ifaid"`
	Data       interface{} `json:"data"`
}

func init() {
	connection.RedisConnect()
	redisdb = connection.GetRedisDb()
}

func SaveToRedis(fileName string, IfaId int, InvestorId int, Data []byte) {
	now := time.Now()
	//key := fmt.Sprintf("%s_%s_%d", fileName, now.Format(YYYYMMDD), InvestorId)
	// err := redisdb.Set(ctx, key, data, 0).Err()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	key := fmt.Sprintf("%s_%d_%s:%d", fileName, IfaId, now.Format(config.YYYYMMDD), InvestorId)
	//err := redisdb.HMSet(ctx, key, "investorid", InvestorId, "ifaid", IfaId, "data", Data).Err()
	err := redisdb.HMSet(ctx, key, map[string]interface{}{
		"investorid": InvestorId, "ifaid": IfaId, "data": Data,
	}).Err()
	if err != nil {
		log.Fatal(err)
	} else {
		//ttl, err := strconv.Atoi(connection.GoDotEnvVariable("REDIS_TTL"))
		if err != nil {
			log.Fatal("Expiry setting error", err)
		}
		_, err = redisdb.Expire(ctx, key, 168*time.Hour).Result()
		if err != nil {
			log.Fatal("Expiry setting error", err)
		} else {
			fmt.Println("Data saved to redis ", InvestorId)
		}
	}

}

func GetAllData(fileName string, IfaId int) []interface{} {
	var response []interface{}
	now := time.Now()
	Hashkey := fmt.Sprintf("%s_%d_%s:*", fileName, IfaId, now.Format(config.YYYYMMDD))
	result, err := redisdb.Keys(ctx, Hashkey).Result()
	if err != nil {
		log.Fatal(err)
	}

	for _, key := range result {
		vals, err := redisdb.HMGet(ctx, key, "data").Result()
		if err != nil {
			log.Fatal(err)
		}
		var data ApiResponse
		if err := json.Unmarshal([]byte(fmt.Sprintf("%v", vals...)), &data); err != nil {
			log.Fatal(err)
		}
		response = append(response, data)
	}
	return response
}

func InitialLogginToRedis(fileType string, IfaId int, Count int64) {
	now := time.Now()
	key := fmt.Sprintf("%s_Initial_%d_%s", fileType, IfaId, now.Format(config.YYYYMMDD))

	err := redisdb.HMSet(ctx, key, map[string]interface{}{
		"IfaId":   IfaId,
		"Total":   Count,
		"Error":   0,
		"Success": 0,
	}).Err()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Data saved and expiry set to redis  for initial logging ")
	}
}

func UpdateLogginToRedis(IfaId int, fileType string, Count int) {
	now := time.Now()
	key := fmt.Sprintf("%s_Initial_%d_%s", fileType, IfaId, now.Format(config.YYYYMMDD))
	vals, err := redisdb.HGetAll(ctx, key).Result()
	if err != nil {
		log.Fatal(err)
	}

	LoggingDetails := vals
	total, err := strconv.Atoi(LoggingDetails["Total"])
	if err != nil {
		log.Fatal("Conversion failed")
	}
	LoggingDetails["Success"] = fmt.Sprint(Count)
	LoggingDetails["Error"] = fmt.Sprint(total - Count)
	redisdb.HSet(ctx, key, LoggingDetails)
}

func GetInvestorData(fileName string, IfaId int, InvestorId int) interface{} {
	now := time.Now()
	return _getInvestorData(fileName, IfaId, InvestorId, now)
}

func GetFailedData(FileType string, IfaId int) []string {
	now := time.Now()
	key := fmt.Sprintf("%s_Failed_%d_%s", FileType, IfaId, now.Format(config.YYYYMMDD))
	response, err := redisdb.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		utils.Error.Println(err)
	}
	return response
}

func SaveFailedData(FileType string, IfaId int, InvestorId int) {
	now := time.Now()
	key := fmt.Sprintf("%s_Failed_%d_%s", FileType, IfaId, now.Format(config.YYYYMMDD))
	if err := redisdb.LPush(ctx, key, InvestorId).Err(); err != nil {
		utils.Error.Println(err)
	}
}

func GetPreviousData(FileName string, IfaId int, investorId int) interface{} {
	//	now := time.Now().AddDate(0, 0, -1)
	now := time.Now()
	return _getInvestorData(FileName, IfaId, investorId, now)
}

func _getInvestorData(fileName string, IfaId int, InvestorId int, now time.Time) interface{} {
	var eresponse interface{}
	Hashkey := fmt.Sprintf("%s_%d_%s:%d", fileName, IfaId, now.Format(config.YYYYMMDD), InvestorId)
	response, err := redisdb.HMGet(ctx, Hashkey, "data").Result()
	if err != nil {
		utils.Error.Println(err)
		return eresponse
	} else {
		var data ApiResponse
		if err := json.Unmarshal([]byte(fmt.Sprintf("%v", response...)), &data); err != nil {
			utils.Error.Println(err)
			return eresponse
		}
		return data
	}
}
