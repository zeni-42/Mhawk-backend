package database

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func ConnectRedis() {
	ctx, cancle := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancle()

	RDB = redis.NewClient(&redis.Options{
		Addr:	  "localhost:6379",
        Password: "", 
        DB:		  0,
        Protocol: 2,
	})

	pong, err := RDB.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis failed")
	}

	log.Println("REDIS", pong)
}

func DisconnectRedis() {
	if RDB != nil {
		if err := RDB.Close(); err != nil {
			log.Printf("Redis failed to disconnect")
			return
		}
		log.Printf("REDIS DISCONNECTED")
	}
}

func GetRedisPing() bool {
	if RDB != nil {
		_, err := RDB.Ping(context.Background()).Result()
		if err != nil {
			return err == nil
		}
		return true
	}
	return false
}