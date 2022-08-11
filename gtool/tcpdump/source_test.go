package tcpdump

import (
	"context"
	"os"
	"testing"

	"github.com/anthony-dong/go-sdk/commons/tcpdump"
)

func TestNewConsulSource(t *testing.T) {
	cfg := tcpdump.NewDefaultConfig()
	cfg.Dump = true
	ctx := tcpdump.NewCtx(context.Background(), cfg)
	ctx.AddDecoder("http", tcpdump.NewHTTP1Decoder())
	ctx.AddDecoder("thrift", tcpdump.NewThriftDecoder())
	testFile := func(file string) {
		open, err := os.Open(readFile(file))
		if err != nil {
			t.Fatal(err)
		}
		defer open.Close()
		source := NewConsulSource(open, NewDecodeOptions())
		for elem := range source.Packets() {
			packet := debugPacket(elem, ctx)
			ctx.HandlerPacket(packet)
			elem.(WaitPacket).Notify()
		}
	}
	testFile("http1.1.pcap.console")
	testFile("thrift.pcap.console")
	testFile("thrift_ttheader.pcap.console")
}
