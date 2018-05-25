package scheduler

import (
	"sync"

	"github.com/go-redis/redis"
)

var (
	setRedisOnce sync.Once
	setRedisErr  error
	redisClient  *redis.Client
)

func SetRedis(redisURL string) error {
	setRedisOnce.Do(func() {
		var opts *redis.Options
		opts, setRedisErr = redis.ParseURL(redisURL)
		if setRedisErr != nil {
			return
		}

		redisClient = redis.NewClient(opts)
		setRedisErr = redisClient.Ping().Err()
	})

	return setRedisErr
}
