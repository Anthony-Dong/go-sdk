package hexo

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/anthony-dong/go-sdk/commons"
	"github.com/anthony-dong/go-sdk/gotool/config"
	git "github.com/sabhiram/go-gitignore"
	"github.com/stretchr/testify/assert"
)

func TestGetAllPage(t *testing.T) {
	config.SetConfigFile("gotool.yaml")
	list, err := GetAllPage("test", []string{})
	if err != nil {
		t.Fatal(err)
	}
	for _, elem := range list {
		t.Log(elem)
	}

}
func TestLines(t *testing.T) {
	t.Run("绝对路径", func(t *testing.T) {
		lines := git.CompileIgnoreLines("/bin")
		assert.Equal(t, lines.MatchesPath("bin/tool"), true)
		assert.Equal(t, lines.MatchesPath("data/bin/tool"), false)
	})
	t.Run("相对路径", func(t *testing.T) {
		lines := git.CompileIgnoreLines(".git")
		assert.Equal(t, lines.MatchesPath("/.git/tool"), true)
		assert.Equal(t, lines.MatchesPath("/data/.git/tool"), true)
	})
}

func TestCheckFileCanHexoPre(t *testing.T) {
	assert.Equal(t, CheckFileCanHexoPre("test/hexo.md"), true)
	assert.Equal(t, CheckFileCanHexoPre("test/not_hexo.md"), false)
}

func TestCheckFileCanHexo(t *testing.T) {
	result, err := CheckFileCanHexo("test/hexo.md", "")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(commons.ToPrettyJsonString(result))
}

func TestRun(t *testing.T) {
	config.SetConfigFile("gotool.yaml")
	dir := filepath.Clean("test")
	targetDir := filepath.Clean("test/post")
	firmCode := []string{"baidu", "ali"}
	if err := buildHexo(context.Background(), dir, targetDir, firmCode, nil); err != nil {
		t.Fatal(err)
	}
}
