package clients

import (
	"github.com/bsm/redislock"
	"time"
)

var lock *redislock.Client

func InitLock() {
	lock = redislock.New(client)
}

func GetLock(key string, duration time.Duration) (*redislock.Lock, error) {
	lock, err := lock.Obtain(key, duration, nil)
	if err != nil {
		return nil, err
	}
	return lock, nil
}

func ReleaseLock(lock *redislock.Lock) error {
	return lock.Release()
}
