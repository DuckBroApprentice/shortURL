package redisConn

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// 大概用不到，順便練習而已
func Insert(ctx context.Context, key string, value interface{}) (string, error) {
	err := RDB.Set(ctx, key, value, 0).Err()
	if err != nil {
		return "insert fail", err
	}
	return "insert success !", nil
}

func Search(ctx context.Context, key string) (string, error) {
	val, err := RDB.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "key does not eixst", err
		}
		return "", err
	}
	return val, nil
}

func Delete(ctx context.Context, keys []string) int64 {
	return RDB.Del(ctx, keys...).Val()
}

func Update(ctx context.Context, key string, value interface{}) {
	err := RDB.Set(ctx, key, value, 0).Err()
	if err != nil {
		fmt.Println(err)
	}
}
