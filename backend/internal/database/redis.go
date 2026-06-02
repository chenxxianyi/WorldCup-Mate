package database

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func InitRedis(addr, password, dbStr string) {
	db, _ := strconv.Atoi(dbStr)
	RDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := RDB.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}
	fmt.Println("Redis connected")
}
