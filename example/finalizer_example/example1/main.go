package main

import (
	"log"
	"runtime"
	"sync"
	"time"
)

func init() {
	log.SetFlags(log.Ltime)
}
func (c *cacheData) Set(key string, v interface{}) {
	c.dataLock.Lock()
	defer c.dataLock.Unlock()
	c.data[key] = v
}

type CacheData struct {
	*cacheData
}

type cacheData struct {
	name     string
	dataLock sync.RWMutex
	data     map[string]interface{}
	reporter func(data *cacheData)

	closeOnce sync.Once
	done      chan struct{}
}

func NewCacheData(name string) *cacheData {
	data := &cacheData{
		name: name,
		data: map[string]interface{}{},
		reporter: func(data *cacheData) {
			log.Println("reporter")
		},
		done: make(chan struct{}, 0),
	}
	return data
}

// NewSafeCacheData 安全的函数
func NewSafeCacheData(name string) *CacheData {
	data := NewCacheData(name)
	data.init()
	result := &CacheData{
		cacheData: data,
	}
	runtime.SetFinalizer(result, (*CacheData).Close)
	return result
}

// init 注册reporter函数，比如上报一些缓存的信息
func (c *cacheData) init() {
	go func() {
		c.reporter(c)
		t := time.NewTicker(time.Second)
		for {
			select {
			case <-c.done:
				t.Stop()
				return
			case <-t.C:
				c.reporter(c)
			}
		}
	}()
}

// Close 函数主要是防止goroutine泄漏
func (c *cacheData) Close() {
	c.closeOnce.Do(func() {
		close(c.done)
	})
}

func BizFunc() {
	cache := NewSafeCacheData("test")

	cache.Set("k1", "v1")

	// biz ....
	// 但是忘记关闭cache了，或者等等的没有close，导致G泄漏
}

func main() {
	BizFunc()
	for x := 0; x < 10; x++ {
		runtime.GC()
		log.Println("runtime.GC")
		time.Sleep(time.Second)
	}
}
