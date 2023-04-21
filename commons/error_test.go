package commons

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestName(t *testing.T) {
	cause := errors.New("whoops")
	err := errors.WithStack(cause)
	fmt.Printf("%+v", err)
}

func TestErrorBUG(t *testing.T) {
	if err := test1(); err != nil {
		t.Fatal(err)
	}
	if err := test3(); err != nil {
		t.Logf("test3() return err: %v, type: %T", err, err)
	} else {
		t.Fatal(`must error`)
	}

	if err := test4(""); err != nil {
		t.Logf("test4(\"\") return err: %v\n", err)
	} else {
		t.Fatal("must error")
	}

	if err := test4("1"); err != nil {
		t.Fatal(err)
	}
}

func test4(v string) error {
	if v == "" {
		return &Errno{V: "empty"}
	}
	return nil
}

func test3() error {
	return test2()
}
func test2() error {
	return test1()
}
func test1() *Errno {
	return nil
}

type Errno struct {
	V string
}

func (r *Errno) Error() string {
	return r.V
}

func TestChannel(t *testing.T) {
	ints := make(chan int, 0)
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer func() {
			close(ints)
			wg.Done()
		}()
		for x := 0; x < 10; x++ {
			ints <- x
		}
	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case d, isDone := <-ints:
				fmt.Printf("1: receive: %d. isDone: %v\n", d, !isDone)
				if !isDone {
					return
				}
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case d, isDone := <-ints:
				fmt.Printf("2: receive: %d. isDone: %v\n", d, !isDone)
				if !isDone {
					return
				}
			}
		}
	}()

	time.Sleep(time.Second * 10)
}
