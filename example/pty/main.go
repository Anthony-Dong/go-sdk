package main

import (
	"io"
	"os"
	"time"

	"github.com/kr/pty"
)

func main() {
	_, tty, err := pty.Open()
	if err != nil {
		panic(err)
	}
	go func() {
		io.Copy(tty, os.Stdin)
	}()
	go func() {
		io.Copy(os.Stdout, tty)
	}()
	time.Sleep(time.Second * 10000)
}
