package scheduler

import (
	"github.com/go-redis/redis"
)

func PushURL(url string) error {
	ok, err := redisClient.SIsMember("dirty_urls", url).Result()
	if err != redis.Nil && err != nil {
		return err
	}
	if ok {
		return nil
	}

	return redisClient.SAdd("urls", url).Err()
}

func PopURL() (url string, isDrained bool, err error) {
	url, err = redisClient.SPop("urls").Result()
	switch {
	case err == redis.Nil || url == "":
		isDrained = true
		return
	case err != nil:
		return
	}

	err = redisClient.SAdd("dirty_urls", url).Err()

	return
}

func PushHTML(html string) error {
	return redisClient.SAdd("htmls", html).Err()
}

func PopHTML() (html string, isDrained bool, err error) {
	html, err = redisClient.SPop("htmls").Result()
	if err == redis.Nil || html == "" {
		isDrained = true
	}

	return
}
