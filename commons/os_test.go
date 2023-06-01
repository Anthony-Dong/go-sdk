package commons

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetGoPath(t *testing.T) {
	data := GetGoPath()
	t.Log(data == "/Users/bytedance/go")

	dir, _ := os.UserHomeDir()
	testPath := filepath.Join(dir, "go")
	t.Log(testPath == "/Users/bytedance/go")
}
