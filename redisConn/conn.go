package redisConn

import (
	"log"
	"shorten/app"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func RedisDB() *redis.Client {
	log.Println("ready connect to redis")
	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	pong, err := RDB.Ping(app.Ctx).Result()
	if err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
	}
	log.Println("Redis connected successfully:", pong)
	// log.Println("redis connect")
	return RDB
}
