package hexo

import (
	"context"
	"testing"
)

func Test_markdownCommand_Run(t *testing.T) {
	cfg := markdownCommand{
		Dir:          "test_readme",
		TemplateFile: "test_readme/README.md.tmpl",
		Ignore:       []string{"/README.md.tmpl", "/README.md"},
	}
	if err := cfg.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
}
