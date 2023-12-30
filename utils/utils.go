package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     GetConfig("REDIS_HOST"),
		Username: GetConfig("REDIS_LOGIN"),
		Password: GetConfig("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	return client
}

func IsAllowedUser(userID int64) bool {
	id := fmt.Sprintf("%d", userID)
	allowedUserIDs := strings.Split(GetConfig("ALLOWED_TELEGRAM_USER_IDS"), ",")

	for _, allowedUserID := range allowedUserIDs {
		if allowedUserID == id {
			return true
		}
	}

	return false
}

func GetConfig(key string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return ""
}
