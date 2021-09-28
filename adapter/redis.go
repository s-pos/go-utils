package adapter

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var (
	// redisdb is variable
	redisdb *redis.Client
)

// LoadRedis is load connectivity to redis database
func LoadRedis() {
	dbnum, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		panic(err)
	}

	addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	client := &redis.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASS"),
		DB:       dbnum,
	}

	redisdb = redis.NewClient(client)
}

// GetClientRedis is for query to redis database
func GetClientRedis() *redis.Client {
	return redisdb
}
