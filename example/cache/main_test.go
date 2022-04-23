package main

import (
	"context"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
)

type localCache struct {
	sync.Map
}

func (l *localCache) Get(key string) (value interface{}, isExist bool) {
	return l.Load(key)
}

func (l *localCache) Set(key string, value interface{}, opts ...CacheOption) {
	l.Store(key, value)
}

type goCache struct {
	*cache.Cache
}

func (l goCache) Set(key string, value interface{}, opts ...CacheOption) {
	l.SetDefault(key, value)
}

func TestNewCached(t *testing.T) {
	cached := NewWrapperCache("test", goCache{
		Cache: cache.New(time.Second*10, time.Second*30),
	})
	//cached := NewWrapperCache("test", &localCache{})
	ctx := context.Background()
	wg := sync.WaitGroup{}
	var (
		loadTime uint64 = 0
		currG           = 20
	)
	wg.Add(currG)
	for x := 0; x < currG; x++ {
		go func(x int) {
			defer wg.Done()
			for y := 0; y < 200000; y++ {
				key := y % 10
				result, err := cached.GetData(ctx, strconv.Itoa(key), func(ctx context.Context) (interface{}, error) {
					atomic.AddUint64(&loadTime, 1)
					t.Logf("load key: %s, y: %d, x: %d\n", strconv.Itoa(key), y, x)
					return int(key), nil
				})
				if err != nil {
					t.Fatal(err)
				}
				if result.(int) != key {
					t.Fatal("data find err")
				}
			}
		}(x)
	}
	wg.Wait()
	for x := 0; x < 10; x++ {
		result, _ := cached.GetData(ctx, strconv.Itoa(x), nil)
		t.Log(result)
		assert.Equal(t, result.(int), int(x))
	}

	assert.Equal(t, int(loadTime), int(10))
}
