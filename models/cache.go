package models

import (
	"context"
	"time"
)

func CacheSet(key string, value interface{}, exptime time.Duration) error {
	err := RedisCli.Set(context.Background(), key, value, exptime).Err()
	if err != nil {
		panic(err)
	}
	return err
}
