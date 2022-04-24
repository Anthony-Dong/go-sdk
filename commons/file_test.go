package commons

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFileRelativePath(t *testing.T) {
	getResult := func(v1, v2 string) string {
		path, err := GetFileRelativePath(v1, v2)
		if err != nil {
			t.Fatal(err)
		}
		return path
	}
	assert.Equal(t, getResult("/data/log/test/a.log", "/data"), "log/test/a.log")

	t.Log(filepath.Rel("/data", "/data/log/test/a.log"))
	rel, _ := filepath.Rel("/data", "/data/log/test/a.log")
	t.Log(filepath.Join("/data", rel))
}

func TestGetGoProjectDir(t *testing.T) {
	t.Log(GetGoProjectDir())
	curDir, _ := filepath.Abs(".")
	t.Log(curDir)
}

func TestGetCmdName(t *testing.T) {
	t.Log(GetCmdName())
}

func TestCheckStdInFromPiped(t *testing.T) {
	t.Log(CheckStdInFromPiped())
}
