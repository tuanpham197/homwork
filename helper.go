package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"time"
)

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func ConvertJSonToStruct[T any](data string, entity *T) error {
	errDecode := json.Unmarshal([]byte(data), &entity)
	if errDecode != nil {
		return errors.New("can't convert json string to struct")
	}

	return nil
}

func GetDatRedis[T any](ctx context.Context, key string, data *T, redisClient *redis.Client) error {
	dataRedis := redisClient.Get(ctx, key)
	jsonData, _ := dataRedis.Result()
	errConvert := ConvertJSonToStruct(jsonData, &data)
	if errConvert != nil {
		return errConvert
	}
	return nil
}

func SetDataToRedis(ctx context.Context, data interface{}, key string, expireTime time.Duration, redisClient *redis.Client) error {
	dataEncoded, errEncode := json.Marshal(data)
	if errEncode != nil {
		return errEncode
	}

	statusCmd := redisClient.Set(ctx, key, dataEncoded, expireTime)
	if statusCmd.Err() != nil {
		return statusCmd.Err()
	}

	return nil
}
