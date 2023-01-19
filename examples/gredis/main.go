package main

import (
	"context"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	glib "github.com/gogf/gf/v2/database/gredis"
	"github.com/jxo-me/redislock"
	"github.com/jxo-me/redislock/redis/gredis"
)

func main() {
	c := &glib.Config{
		Address: "172.20.203.5:6379",
		Pass:    "aa123456",
		Db:      1,
	}
	ctx := context.Background()

	r, err := glib.New(c)
	if err != nil {
		panic(err)
	}
	defer r.Close(ctx)

	rs := redislock.New(gredis.NewPool(r))
	mutex := rs.NewMutex("test-redisLock")

	if err := mutex.Lock(); err != nil {
		panic(err)
	}
	if _, err := mutex.Unlock(); err != nil {
		panic(err)
	}
}
