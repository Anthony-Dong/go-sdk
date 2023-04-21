package protoc

import (
	"path/filepath"
	"testing"
)

func Test_relativePwdPath(t *testing.T) {
	rel, _ := filepath.Rel("/api/v1", "/api/v1/")
	t.Log(rel)
}

func Test_sortFiles(t *testing.T) {
	input := []string{"a/b/c", "a/a/c", "a/b", "a/c", "b", "a/d", "a/cc", "a", "."}
	sortFiles(input)
	t.Logf("%#v\n", input)
	//assert.Equal(t, input, []string{".", "a", "b", "a/b", "a/c", "a/a/c", "a/b/c"})
}
