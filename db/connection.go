package db

import (
	"fmt"

	"github.com/go-redis/redis"
)

func GetConn() *redis.Client {
	farmerDB := redis.NewClient(&redis.Options{
		Addr:     "redis-15262.c20033.asia-seast2-mz.gcp.cloud.rlrcp.com:15262",
		Password: "EG3M5zcZpzqgfdUZcggxfw4koA7DQdvp", // no password set
		// Addr:     "localhost:6379",
		// Password: "",
		// Addr:     "redis-10715.c301.ap-south-1-1.ec2.cloud.redislabs.com:10715",
		// Password: "jeZTdPCPTWdrrSnbjkQZHwS2m1jO1vyy",
		DB: 0, // use default DB
	})
	return farmerDB
}

func PrepareDB() int64 {
	farmerDB := GetConn()
	res, err := farmerDB.SAdd("users:farmer", "0000000000").Result()
	if err != nil {
		fmt.Println(err, "error occurred")
		return 0
	} else {
		return res
	}

}
