package tcpdump

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/anthony-dong/go-sdk/commons/tcpdump"
)

func TestNewConsulSource(t *testing.T) {
	cfg := tcpdump.NewDefaultConfig()
	cfg.Dump = true
	ctx := tcpdump.NewCtx(context.Background(), cfg)
	ctx.AddDecoder("http", tcpdump.NewHTTP1Decoder())
	ctx.AddDecoder("thrift", tcpdump.NewThriftDecoder())
	testFile := func(file string, want int) {
		open, err := os.Open(readFile(file))
		if err != nil {
			t.Fatal(err)
		}
		defer open.Close()
		source := NewConsulSource(open, NewDecodeOptions())
		num := 0
		for elem := range source.Packets() {
			packet := debugPacket(elem, ctx)
			ctx.HandlerPacket(packet)
			elem.(WaitPacket).Notify()
			num++
		}
		assert.Equal(t, num, want)
	}
	testFile("http1.1.pcap.console", 10)
	testFile("thrift.pcap.console", 2)
	testFile("thrift_ttheader.pcap.console", 22)
}
