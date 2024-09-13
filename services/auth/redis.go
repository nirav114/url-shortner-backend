package auth

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"github.com/nirav114/url-shortner-backend.git/types"
)

var ctx = context.Background()

var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Network:    "tcp",
		Addr:       "localhost:6789",
		DB:         0,
		MaxRetries: 2,
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

	return RedisClient.Set(ctx, email+":user", jsonData, OTPExpiry).Err()
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
