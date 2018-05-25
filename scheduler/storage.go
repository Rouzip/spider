package scheduler

import (
	"context"

	"github.com/go-redis/redis"
)

func PushURL(ctx context.Context, url string) error {
	rdc := redisClient.WithContext(ctx)

	ok, err := rdc.SIsMember("dirty_urls", url).Result()
	if err != redis.Nil && err != nil {
		return err
	}
	if ok {
		return nil
	}

	return rdc.SAdd("urls", url).Err()
}

func PopURL(ctx context.Context) (url string, isDrained bool, err error) {
	rdc := redisClient.WithContext(ctx)

	url, err = rdc.SPop("urls").Result()
	switch {
	case err == redis.Nil || url == "":
		isDrained = true
		return
	case err != nil:
		return
	}

	err = rdc.SAdd("dirty_urls", url).Err()

	return
}

func PushHTML(ctx context.Context, html string) error {
	return redisClient.WithContext(ctx).SAdd("htmls", html).Err()
}

func PopHTML(ctx context.Context) (html string, isDrained bool, err error) {
	html, err = redisClient.WithContext(ctx).SPop("htmls").Result()
	if err == redis.Nil || html == "" {
		isDrained = true
	}

	return
}
