package commons

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

var (
	gopath     string
	gopathOnce sync.Once
)

const (
	gopathName = "GOPATH"
)

func GetGoPath() string {
	gopathOnce.Do(func() {
		// load  env
		gopath = os.Getenv(gopathName)
		if gopath != "" {
			return
		}

		// load go env
		stdOut := bytes.Buffer{}
		command := exec.Command("go", "env", gopathName)
		command.Stdout = &stdOut
		if err := command.Run(); err == nil {
			gopath = strings.TrimSuffix(stdOut.String(), "\n")
		}

		// load default home
		if gopath == "" {
			dir, err := os.UserHomeDir()
			if err == nil {
				gopath = filepath.Join(dir, "go")
			}
		}

		// load default name
		if gopath == "" {
			gopath = "/go"
		}
	})

	return gopath
}
