package shell

import (
	"testing"
)

func TestCmd(t *testing.T) {
	err := Cmd("ls -al")
	if err != nil {
		t.Fatal(err)
	}
}
