# Redislock

Redislock provides a gredis-based distributed mutual exclusion lock implementation for Go

## Installation

Install Redislock using the go get command:

    $ go get github.com/jxo-me/redislock

Two driver implementations will be installed; however, only the one used will be included in your project.

 * [GRedis](github.com/gogf/gf/v2/database/gredis)

See the [examples](examples) folder for usage of each driver.

## Usage

Error handling is simplified to `panic` for shorter example.

```go
package main

import (
	"context"
	glib "github.com/gogf/gf/v2/database/gredis"
	"github.com/jxo-me/redislock"
	"github.com/jxo-me/redislock/redis/gredis"
)

func main() {
	// Create a pool with gredis which is the pool redislock will
	// use while communicating with Redis. This can also be any pool that
	// implements the `redis.Pool` interface.
	c := &glib.Config{
		Address: "127.0.0.1:6379",
	}
	ctx := context.Background()
	r, err := glib.New(c)
	if err != nil {
		panic(err)
	}
	defer r.Close(ctx)

	// Create an instance of redislock to be used to obtain a mutual exclusion
	// lock.
	rs := redislock.New(gredis.NewPool(r))
	mutex := rs.NewMutex("test-redisLock")

	// Obtain a lock for our given mutex. After this is successful, no one else
	// can obtain the same lock (the same mutex name) until we unlock it.
	if err := mutex.Lock(); err != nil {
		panic(err)
	}
	// Do your work that requires the lock.

	// Release the lock so other processes or threads can obtain a lock.
	if _, err := mutex.Unlock(); err != nil {
		panic(err)
	}
}
```

## Contributing

Contributions are welcome.
