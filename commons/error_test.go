package commons

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
)

func TestName(t *testing.T) {
	cause := errors.New("whoops")
	err := errors.WithStack(cause)
	fmt.Printf("%+v", err)
}
