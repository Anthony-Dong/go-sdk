package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/cenkalti/backoff"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	var (
		data  interface{}
		start = time.Now()
		next  = start
	)

	if err := backoff.Retry(func() error {
		if rand.Int()%1000 != 0 { // 模拟异常
			err := fmt.Errorf("find data mod 2 is zero")
			log.Printf("find err, err: %s, spend: %dms, inc: %dms\n", err, time.Now().Sub(start)/time.Millisecond, time.Now().Sub(next)/time.Millisecond)
			next = time.Now()
			return err
		}
		data = "load success"
		return nil
	}, backoff.NewExponentialBackOff()); err != nil {
		panic(err)
	}

	log.Printf("data: %s\n", data)
}

//output
//find err, err: find data mod 2 is zero
//data: load success
