package auth

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nirav114/url-shortner-backend.git/config"
	"github.com/nirav114/url-shortner-backend.git/types"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.EnvConfig.REDIS_HOST,
		Password: "",
		DB:       0,
	})

	pong, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	println("Connected to Redis:", pong)
}

func StoreUserOTPData(email string, user types.UserOTPData) error {
	jsonData, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return RedisClient.Set(ctx, email+":user", jsonData, 10*time.Minute).Err()
}

func RetrieveUserOTPData(email string) (*types.UserOTPData, error) {
	jsonData, err := RedisClient.Get(ctx, email+":user").Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var user types.UserOTPData
	err = json.Unmarshal([]byte(jsonData), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func DeleteUserOTPData(email string) error {
	return RedisClient.Del(ctx, email+":user").Err()
}
