package tcpdump

import (
	"os"
	"testing"

	"github.com/google/gopacket/layers"

	"github.com/anthony-dong/go-sdk/gtool/tcpdump/reassembly"

	"github.com/stretchr/testify/assert"

	"github.com/anthony-dong/go-sdk/commons/tcpdump"
)

func TestNewConsulSource(t *testing.T) {
	assembler := reassembly.NewAssembler(tcpdump.NewDefaultDecoder(false, nil, map[string]tcpdump.Decoder{
		"HTTP":   tcpdump.NewHTTP1Decoder(),
		"Thrift": tcpdump.NewThriftDecoder(),
	}))
	testFile := func(file string, want int) {
		open, err := os.Open(readFile(file))
		if err != nil {
			t.Fatal(err)
		}
		defer open.Close()
		source := NewConsulSource(open, NewDecodeOptions())
		num := 0
		for packet := range source.Packets() {
			tcp := packet.TransportLayer().(*layers.TCP)
			assembler.AssembleWithContext(packet.NetworkLayer().NetworkFlow(), tcp, &reassembly.Context{CaptureInfo: packet.Metadata().CaptureInfo})
			packet.(WaitPacket).Notify()
			num++
		}
		assert.Equal(t, num, want)
	}
	testFile("http1.1.pcap.console", 10)
	testFile("thrift.pcap.console", 2)
	testFile("thrift_ttheader.pcap.console", 22)
}
