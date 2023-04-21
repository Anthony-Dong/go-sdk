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

// CGO_ENABLED=1
func Test_DecodeTCPDump(t *testing.T) {
	ctx := context.Background()
	cfg := NewDefaultConfig()
	initCfg := func() {
		cfg = NewDefaultConfig()
	}
	t.Run("thrift", func(t *testing.T) {
		//cfg.DisableReassembly = true
		if err := run(ctx, readFile("thrift.pcap"), cfg); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("http", func(t *testing.T) {
		initCfg()
		if err := run(context.Background(), readFile("http1.1.pcap"), cfg); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("http chunked", func(t *testing.T) {
		initCfg()
		cfg.Loopback = true
		if err := run(context.Background(), readFile("test.pcap"), cfg); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("stick http", func(t *testing.T) {
		initCfg()
		if err := run(context.Background(), readFile("stick_http1.1.pcap"), cfg); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("thrift_ttheader", func(t *testing.T) {
		initCfg()
		// thrift_ttheader
		if err := run(context.Background(), readFile("thrift_ttheader.pcap"), cfg); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("stick_thrift_ttheader", func(t *testing.T) {
		initCfg()
		// stick thrift_ttheader
		if err := run(ctx, readFile("stick_thrift_ttheader.pcap"), cfg); err != nil {
			t.Fatal(err)
		}
	})
}
