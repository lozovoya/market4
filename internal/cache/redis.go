package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

const TTL = 30

type apiCache struct {
	pool *redis.Pool
}

func NewRedisCache(pool *redis.Pool) Cache {
	return &apiCache{pool: pool}
}

func InitCache(addr string) *redis.Pool {
	myCache := &redis.Pool{
		DialContext: func(ctx context.Context) (redis.Conn, error) {
			return redis.DialURL(addr)
		},
	}
	return myCache
}
func (a *apiCache) ToCache(ctx context.Context, key string, value []byte) (err error) {
	conn, err := a.pool.GetContext(ctx)
	if err != nil {
		return fmt.Errorf("ToCache: %w", err)
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Println(err)
		}
	}()

	_, err = redis.DoWithTimeout(conn, time.Millisecond*100, "SETEX", key, TTL, value)
	if err != nil {
		return fmt.Errorf("ToCache: %w", err)
	}
	return nil
}

func (a *apiCache) FromCache(ctx context.Context, key string) (value []byte, err error) {
	conn, err := a.pool.GetContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("FromCache: %w", err)
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Println(cerr)
		}
	}()

	reply, err := redis.DoWithTimeout(conn, time.Millisecond*100, "GET", key)
	if err != nil {
		return nil, fmt.Errorf("FromCache: %w", err)
	}
	value, err = redis.Bytes(reply, nil)
	if err != nil {
		return nil, err
	}

	return value, nil
}
