package tcpdump

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/anthony-dong/go-sdk/commons"
)

func readFile(file string) string {
	dir := commons.GetGoProjectDir()
	return filepath.Join(dir, "tcpdump/test", file)
}

func Test_run(t *testing.T) {
	t.Run("thrift", func(t *testing.T) {
		if err := run(context.Background(), readFile("thrift.pcap"), Thrift, true); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("http", func(t *testing.T) {
		if err := run(context.Background(), readFile("http1.1.pcap"), HTTP, false); err != nil {
			t.Fatal(err)
		}
	})
}
