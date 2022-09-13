package model

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

//定义一个全局的pool
var Pool *redis.Pool

func InitPool(address string, maxIdle int, maxActive int, idleTimeout time.Duration) {
	Pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", address, redis.DialPassword("123456"))
		},
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: idleTimeout,
	}
}
