package main

import (
	"fmt"
	"log"
	"runtime"
	"time"
)

type object int

func (o object) Ptr() *object {
	return &o
}

var (
	cacheData = make(map[string]*object, 1024)
)

func deleteData(key string) {
	delete(cacheData, key)
}

var (
	globalData *object
)

func setData(key string, v object) {
	data := v.Ptr()
	runtime.SetFinalizer(data, func(data *object) error {
		fmt.Printf("runtime invoke Finalizer data: %d, time: %s\n", *data, time.Now().Format("15:04:05.000"))
		time.Sleep(time.Second)
		//panic("我故意panic的")
		if *data == 1 {
			globalData = data //再次回收！
			fmt.Printf("invoke Finalizer data: %d, data: %d, g_data: %d, time: %s\n", *data, data, globalData, time.Now().Format("15:04:05.000"))
		}
		return fmt.Errorf("我故意Error的")
	})
	cacheData[key] = data
}

func main() {
	defer func() {
		err := recover()
		log.Printf("err: %v\n", err)
	}()
	setData("key1", 1)
	setData("key2", 2)
	setData("key3", 3)

	deleteData("key1")
	deleteData("key2")
	deleteData("key3")

	for x := 0; x < 15; x++ {
		fmt.Println("invoke runtime.GC()")
		runtime.GC()
		time.Sleep(time.Second)

		if x == 6 {
			fmt.Printf("set globalData is nil, g_dir: %d, g: %d\n", globalData, *globalData)
			globalData = nil
		}
	}
}

// output:
// invoke runtime.GC()
//runtime invoke Finalizer data: 2, time: 22:05:46.040
//invoke runtime.GC()
//runtime invoke Finalizer data: 1, time: 22:05:47.044
//invoke runtime.GC()
//invoke Finalizer data: 1, data: 824633827472, g_data: 824633827472, time: 22:05:48.045
//runtime invoke Finalizer data: 3, time: 22:05:48.046
//invoke runtime.GC()
//invoke runtime.GC()
//invoke runtime.GC()
//invoke runtime.GC()
//set globalData is nil, g_dir: 824633827472, g: 1
//invoke runtime.GC()
//invoke runtime.GC()
//invoke runtime.GC()
//invoke runtime.GC()
//invoke runtime.GC()
//invoke runtime.GC()
//invoke runtime.GC()
//invoke runtime.GC()
// 22:06:01 err: <nil>
