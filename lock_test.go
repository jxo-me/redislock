package redislock

import (
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	glib "github.com/gogf/gf/v2/database/gredis"
	"github.com/jxo-me/redislock/redis"
	"github.com/jxo-me/redislock/redis/gredis"
	"testing"
)

type testCase struct {
	poolCount int
	pools     []redis.Pool
}

func makeCases(poolCount int) map[string]*testCase {
	return map[string]*testCase{
		"gredis": {
			poolCount,
			newMockPoolsGredis(poolCount),
		},
	}
}

func TestMutex_Lock(t *testing.T) {
	for k, v := range makeCases(8) {
		t.Run(k, func(t *testing.T) {
			rs := New(v.pools...)

			mutex := rs.NewMutex("test-redislock")
			_ = mutex.Lock()

			assertAcquired(t, v.pools, mutex)
		})
	}
}

func newMockPoolsGredis(n int) []redis.Pool {
	pools := make([]redis.Pool, n)
	c := &glib.Config{
		Address: "test.com:6379",
		Pass:    "aa123456",
	}
	//ctx := context.Background()
	glib.SetConfig(c, "test")
	r := glib.Instance("test")
	//defer r.Close(ctx)
	for i := 0; i < n; i++ {
		pools[i] = gredis.NewPool(r)
	}
	return pools
}
