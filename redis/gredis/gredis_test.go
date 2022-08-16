package gredis

import "github.com/jxo-me/redislock/redis"

var _ redis.Conn = (*conn)(nil)

var _ redis.Pool = (*pool)(nil)
