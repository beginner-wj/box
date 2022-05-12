package redis

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

var rp *RedisPool

type RedisPool struct {
	Pool *redis.Pool
}

func InitRedis(RedisHost, RedisPwd string, database ...int) {
	var err error
	db := 0
	if database != nil {
		db = database[0]
	}
	rp, err = redisDB(RedisHost, RedisPwd, db)
	if err != nil {
		panic(err)
	}
}

func GetPool() redis.Conn {
	return rp.Pool.Get()
}

func redisDB(server, passwd string, db int) (*RedisPool, error) {
	rp, err := myRedisPool(server, passwd, db)
	return rp, err
}

func myRedisPool(server, password string, database ...int) (*RedisPool, error) {
	db := 0
	if database != nil {
		db = database[0]
	}
	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}

			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			if _, err := c.Do("SELECT", db); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return &RedisPool{pool}, nil
}

func SetEx(key string, value string, time int32) {
	c := GetPool()
	defer c.Close()
	if _, err := c.Do("SET", key, value, "EX", time); err != nil {
		fmt.Println("redis set:", err)
	}
}
