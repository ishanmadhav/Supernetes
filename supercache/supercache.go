package supercache

import (
	"context"
	"time"

	"github.com/allegro/bigcache/v3"
)

type SuperCache struct {
	Cache *bigcache.BigCache
}

func NewSuperCache() (*SuperCache, error) {
	cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(100*time.Minute))
	if err != nil {
		return nil, err
	}
	return &SuperCache{Cache: cache}, nil
}

func (s *SuperCache) Set(key string, value []byte) error {
	return s.Cache.Set(key, value)
}

func (s *SuperCache) Get(key string) (string, error) {
	bytes, err := s.Cache.Get(key)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (s *SuperCache) Exists(key string) bool {
	_, err := s.Cache.Get(key)
	if err != nil {
		return false
	} else {
		return true
	}
	return true
}

func (s *SuperCache) Delete(key string) error {
	return s.Cache.Delete(key)
}

func (s *SuperCache) Reset() error {
	return s.Cache.Reset()
}

func (s *SuperCache) Close() error {
	return s.Cache.Close()
}

func (s *SuperCache) Len() int {
	return s.Cache.Len()
}
