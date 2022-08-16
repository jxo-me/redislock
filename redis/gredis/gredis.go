package gredis

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gredis"
	"strings"
	"time"

	"github.com/jxo-me/redislock/redis"
)

type pool struct {
	delegate *gredis.Redis
}

func (p *pool) Get(ctx context.Context) (redis.Conn, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	c, err := p.delegate.Conn(ctx)
	if err != nil {
		return nil, err
	}
	return &conn{c, ctx}, nil
}

// NewPool returns a Goredis-based pool implementation.
func NewPool(delegate *gredis.Redis) redis.Pool {
	return &pool{delegate}
}

type conn struct {
	delegate *gredis.RedisConn
	ctx      context.Context
}

func (c *conn) Get(name string) (string, error) {
	value, err := c.delegate.Do(c.ctx, "GET", name)
	return value.String(), noErrNil(err)
}

func (c *conn) Set(name string, value string) (bool, error) {
	reply, err := c.delegate.Do(c.ctx, "SET", name, value)
	return reply.String() == "OK", err
}

func (c *conn) SetNX(name string, value string, expiry time.Duration) (bool, error) {
	reply, err := c.delegate.Do(c.ctx, "SET", name, value, "NX", "PX", expiry*time.Millisecond)
	return reply.Bool(), err
}

func (c *conn) PTTL(name string) (time.Duration, error) {
	expiry, err := c.delegate.Do(c.ctx, "PTTL", name)
	fmt.Println("expiry:", expiry)
	return time.Duration(expiry.Int64()) * time.Millisecond, noErrNil(err)
}

func (c *conn) Eval(script *redis.Script, keysAndArgs ...interface{}) (interface{}, error) {
	v, err := c.delegate.Do(c.ctx, "EVALSHA", args(script, script.Hash, keysAndArgs)...)
	if err != nil {
		if strings.Contains(err.Error(), "NOSCRIPT ") {
			v, err = c.delegate.Do(c.ctx, "EVAL", args(script, script.Src, keysAndArgs)...)
		}
	}
	return v, noErrNil(err)
}

func (c *conn) Close() error {
	// Not needed for this library
	if c.delegate != nil {
		err := c.delegate.Close(c.ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func noErrNil(err error) error {
	if err != nil {
		return err
	}
	return nil
}

func args(script *redis.Script, spec string, keysAndArgs []interface{}) []interface{} {
	var args []interface{}
	if script.KeyCount < 0 {
		args = make([]interface{}, 1+len(keysAndArgs))
		args[0] = spec
		copy(args[1:], keysAndArgs)
	} else {
		args = make([]interface{}, 2+len(keysAndArgs))
		args[0] = spec
		args[1] = script.KeyCount
		copy(args[2:], keysAndArgs)
	}
	return args
}
