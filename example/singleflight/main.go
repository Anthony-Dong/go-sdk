package main

import (
	"fmt"
	"sync"

	"github.com/golang/groupcache/singleflight"
)

var (
	cache sync.Map
	sf    = singleflight.Group{}
)

func main() {
	key := "k1" // 假如现在有100个并发请求访问 k1
	wg := sync.WaitGroup{}
	wg.Add(100)
	for x := 0; x < 100; x++ {
		go func() {
			defer wg.Done()
			loadKey(key)
		}()
	}
	wg.Wait()
	fmt.Printf("result key: %s\n", loadKey(key))
}
func loadKey(key string) (v string) {
	if data, ok := cache.Load(key); ok {
		return data.(string)
	}
	data, err := sf.Do(key, func() (interface{}, error) {
		if data, ok := cache.Load(key); ok { // 双重检测
			return data.(string), nil
		}
		data := "data" + "|" + key
		fmt.Printf("load and set success, data: %s\n", data)
		cache.Store(key, data)
		return data, nil
	})
	if err != nil {
		// todo handler
		panic(err)
	}
	return data.(string)
}
