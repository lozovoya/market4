package cache

import "context"

type Cache interface {
	ToCache(ctx context.Context, key string, value []byte) (err error)
	FromCache(ctx context.Context, key string) (value []byte, err error)
}
