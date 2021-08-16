package cache

import (
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

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
			log.Printf("ToCache: %w", err)
		}
	}()

	_, err = redis.DoWithTimeout(conn, time.Millisecond*100, "SET", key, value)
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
			log.Printf("FromCache: %w", err)
		}
	}()

	reply, err := redis.DoWithTimeout(conn, time.Millisecond*100, "GET", key)
	if err != nil {
		return nil, fmt.Errorf("FromCache: %w", err)
	}
	value, _ = redis.Bytes(reply, err)
	//todo обработать ошибку

	return value, nil
}
