package redis

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

var Pool *redis.Pool

func init() {
	Pool = &redis.Pool{
		MaxIdle:     100,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "114.215.146.129:6379", redis.DialPassword("123456"))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
