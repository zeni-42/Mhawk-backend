package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeni-42/Mhawk/internal/models"
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

func SetUserDataInRedis(k, v string) error {
	if RDB == nil {
		log.Println("Redis client is empty")
		return fmt.Errorf("client is empty")
	}

	ctx := context.Background()
	err := RDB.Set(ctx, k, v, 5*time.Minute).Err()
	if err != nil {
		log.Printf("Failed to set Data in Redis %v", err)
		return err
	}

	return nil
}

func GetUserDataFromRedis(k string) (models.User, error) {
	if RDB == nil {
		log.Println("Redis client is empty")
		return models.User{}, fmt.Errorf("client is empty")
	}

	ctx := context.Background()
	data, err := RDB.Get(ctx, k).Result()
	if err != nil {
		return models.User{}, err
	}

	var jsonData models.User

	if err := json.Unmarshal([]byte(data), &jsonData); err != nil {
		log.Println("Unmarshaling failed")
		return models.User{}, err
	}

	return jsonData, nil
}