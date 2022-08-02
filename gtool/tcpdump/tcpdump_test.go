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
	ctx := context.Background()
	//t.Run("thrift", func(t *testing.T) {
	//	if err := run(context.Background(), readFile("thrift.pcap"), Thrift, false); err != nil {
	//		t.Fatal(err)
	//	}
	//})
	//t.Run("http", func(t *testing.T) {
	//	if err := run(context.Background(), readFile("http1.1.pcap"), HTTP, false); err != nil {
	//		t.Fatal(err)
	//	}
	//})
	//t.Run("thrift_ttheader", func(t *testing.T) {
	//	// thrift_ttheader
	//	if err := run(context.Background(), readFile("thrift_ttheader.pcap"), Thrift, false); err != nil {
	//		t.Fatal(err)
	//	}
	//})
	t.Run("thrift_ttheader2", func(t *testing.T) {
		// thrift_ttheader
		if err := run(ctx, readFile("out.pcap"), false, false); err != nil {
			t.Fatal(err)
		}
	})
}
