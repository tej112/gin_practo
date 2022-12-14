package db

import (
	"encoding/json"
	"fmt"
	"log"
	"main/models"

	"github.com/go-redis/redis"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/nitishm/go-rejson"
)

var farmerDB = GetConn()
var rh = rejson.NewReJSONHandler()

func Exists(key string) *redis.BoolCmd {
	return farmerDB.SIsMember("users:farmer", key)
}

func HashOfFarmer(key string) *redis.StringStringMapCmd {
	return farmerDB.HGetAll(key)
}

func ValueForField(key string, field string) *redis.StringCmd {
	return farmerDB.HGet(key, field)
}

func GetFarmer(key string) interface{} {

	rh.SetGoRedisClient(farmerDB)
	studentJSON, err := redigo.Bytes(rh.JSONGet(":farmer_data_model.Farmer:"+key, "."))
	if err != nil {
		fmt.Println(err.Error())
		// log.Fatalf("Failed to JSONGet")
		return false
	}
	readStudent := models.Farmer{}
	err = json.Unmarshal(studentJSON, &readStudent)
	if err != nil {
		log.Fatalf("Failed to JSON Unmarshal")
	}
	return readStudent
}

func CreateNewFarmer(farmer models.Farmer) {
	rh.SetGoRedisClient(farmerDB)
	res, err := rh.JSONSet(":farmer_data_model.Farmer:"+farmer.Pk, ".", farmer)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(res)
	farmerDB.HSet(farmer.Contact_num, "pk", farmer.Pk)
	farmerDB.HSet(farmer.Contact_num, "language_preference", farmer.Language_preference)

	farmerDB.SAdd("users:farmer", farmer.Contact_num)
}
