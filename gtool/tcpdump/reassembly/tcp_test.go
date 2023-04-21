package reassembly

import (
	"net"
	"testing"
)

func TestName(t *testing.T) {
	//file := `/Users/bytedance/go/src/github.com/anthony-dong/go-sdk/gtool/tcpdump/test/http1.1.pcap`
	//assembler := NewAssembler(true, []commons_tcpdump.Decoder{commons_tcpdump.NewThriftDecoder()})
	//source, err := tcpdump.NewFileSource(file, tcpdump.NewDecodeOptions())
	//if err != nil {
	//	t.Fatal(err)
	//}
	//for packet := range source.Packets() {
	//	tcp := packet.TransportLayer().(*layers.TCP)
	//	assembler.AssembleWithContext(packet.NetworkLayer().NetworkFlow(), tcp, &Context{CaptureInfo: packet.Metadata().CaptureInfo})
	//}
}

func TestIp(t *testing.T) {
	t.Log(net.ParseIP("::1").String())
}
