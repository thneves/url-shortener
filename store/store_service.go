package store

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type StorageService struct {
	redisClient *redis.Client
}

var (
	storeService = &StorageService{}
	ctx          = context.Background()
)

// Note that in a real world usage, the cache duration shouldn't have
// an expiration time, an LRU policy config should be set where the
// values that are retrieved less often are purged automatically from
// the cache and stored back in RDBMS whenever the cache is full

const CacheDuration = 6 * time.Hour

func NewStoreService() *StorageService {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := redisClient.Ping(ctx).Result()

	if err != nil {
		panic(fmt.Sprintf("Error init Redis: %v", err))
	}

	fmt.Printf("\nRedis started succesfully: pong message = {%s}", pong)
	storeService.redisClient = redisClient
	return storeService
}

// We want to be able to save the mapping between the originalUrl and the generated shortUrl url

func SaveUrlMapping(shortUrl string, originalUrl string, userId string) {
	err := storeService.redisClient.Set(ctx, shortUrl, originalUrl, CacheDuration).Err()

	if err != nil {
		panic(fmt.Sprintf("Failed saving url Key | Error: %v - shortUrl: %s - originalUrl: %s\n", err, shortUrl, originalUrl))
	}
}

// We should be able to retrieve the initial long URL once the short is provbided.

func RetrieveInitialUrl(shortUrl string) string {
	result, err := storeService.redisClient.Get(ctx, shortUrl).Result()

	if err != nil {
		panic(fmt.Sprintf("Failed retrieving Original URL | Error: %v - shortUrl: %s/n", err, shortUrl))
	}

	return result
}
