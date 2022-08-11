package bufutils

import (
	"io"
	"time"
)

type WReader interface {
	io.Reader
	io.Writer
	Peek(int) ([]byte, error)
}

type TimeoutReader interface {
	SetReadDeadline(t time.Time) error
}
