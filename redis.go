package utee

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

func CreateRedisPool(size int, server, auth string) *redis.Pool {
	pool := &redis.Pool{
		MaxIdle:     size,
		MaxActive:   size,
		Wait:        true,
		IdleTimeout: 4 * time.Minute,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if auth != "" {
				if _, err := c.Do("AUTH", auth); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			//test only if the connection is not used within last 30 seconds
			if time.Now().After(t.Add(10 * time.Second)) {
				_, err := c.Do("PING")
				return err
			} else {
				return nil
			}
		},
	}
	return pool
}
