package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	"golang.org/x/sync/singleflight"
)

var (
	localCacheCallbackIsNil = fmt.Errorf("cache callback func is nil")
)

type CacheOption interface {
}
type Cache interface {
	Get(key string) (value interface{}, isExist bool)
	Set(key string, value interface{}, opts ...CacheOption)
}
type WrapperCache interface {
	GetData(ctx context.Context, key string, callback func(ctx context.Context) (interface{}, error)) (v interface{}, err error)
}

type wrapperCache struct {
	name           string
	cache          Cache
	singleflight   singleflight.Group
	retrySleepTime time.Duration
	retryNum       uint64
}

func NewWrapperCache(name string, cache Cache) WrapperCache {
	return &wrapperCache{
		name:           name,
		cache:          cache,
		retryNum:       3,
		retrySleepTime: time.Millisecond * 10,
	}
}

// emitHitCachedMetric 计算缓存命中率
func (c *wrapperCache) emitHitCachedMetric(hit bool) {

}
func (c *wrapperCache) GetData(ctx context.Context, key string, callback func(ctx context.Context) (interface{}, error)) (v interface{}, err error) {
	if result, isExist := c.cache.Get(key); isExist {
		c.emitHitCachedMetric(true)
		return result, nil
	}
	if callback == nil {
		return nil, localCacheCallbackIsNil
	}
	c.emitHitCachedMetric(false)

	result, err, _ := c.singleflight.Do(key, func() (interface{}, error) {
		// 双重检测，防止singleflight 锁的key失效
		if result, isExist := c.cache.Get(key); isExist {
			return result, nil
		}
		var callBackData interface{}
		if err := backoff.Retry(func() error {
			if data, err := callback(ctx); err != nil {
				return err
			} else {
				callBackData = data
				return nil
			}
		}, backoff.WithMaxRetries(backoff.NewConstantBackOff(c.retrySleepTime), c.retryNum)); err != nil {
			// todo add log
			return nil, err
		}
		c.cache.Set(key, callBackData)
		return callBackData, nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}
