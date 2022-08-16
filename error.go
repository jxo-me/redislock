package redislock

import "errors"

// ErrFailed is the error resulting if Lock fails to acquire the lock after
// exhausting all retries.
var ErrFailed = errors.New("redisLock: failed to acquire lock")

// ErrExtendFailed is the error resulting if Lock fails to extend the
// lock.
var ErrExtendFailed = errors.New("redisLock: failed to extend lock")
