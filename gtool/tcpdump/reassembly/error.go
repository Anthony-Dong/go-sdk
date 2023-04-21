package reassembly

import (
	"fmt"

	"github.com/fatih/color"
)

type tcpError struct {
	Type string
	Name string
}

func newTCPError(Type string, s string, v ...interface{}) error {
	return tcpError{
		Type: Type,
		Name: fmt.Sprintf(s, v...),
	}
}
func (t tcpError) Error() string {
	return color.RedString(fmt.Sprintf(`[%s] %s`, t.Type, t.Name))
}
