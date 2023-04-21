package tcpdump

import (
	"context"
	"os"

	"github.com/fatih/color"
	"github.com/google/gopacket/layers"
	"github.com/spf13/cobra"

	"github.com/anthony-dong/go-sdk/commons"
	"github.com/anthony-dong/go-sdk/commons/codec"
	"github.com/anthony-dong/go-sdk/commons/tcpdump"
	"github.com/anthony-dong/go-sdk/gtool/tcpdump/reassembly"
	"github.com/anthony-dong/go-sdk/gtool/tcpdump/utils"
)

func NewCmd() (*cobra.Command, error) {
	var (
		cfg      = NewDefaultConfig()
		filename string
	)
	cmd := &cobra.Command{
		Use:   `tcpdump [-r file] [-v] [-X] [--max dump size]`,
		Short: `decode tcpdump file`,
		Long:  `decode tcpdump file, help doc: https://github.com/Anthony-Dong/go-sdk/tree/master/gtool/tcpdump`,
		Example: `  sudo tcpdump -i eth0  -n  -l -X | gtool tcpdump

Help Doc:
  - https://www.tcpdump.org/manpages/pcap-filter.7.html
  - https://www.tcpdump.org/manpages/tcpdump.1.html`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context(), filename, cfg)
		},
	}
	cmd.Flags().BoolVar(&cfg.Verbose, "verbose", false, "Enable Verbose.")
	cmd.Flags().BoolVar(&cfg.DisableReassembly, "disable_reassembly", false, "Disable tcp Reassembly.")
	cmd.Flags().BoolVar(&cfg.Loopback, "loopback", false, "The NIC type for packet capture is loopback.")
	cmd.Flags().StringVarP(&filename, "file", "r", "", "The packets file, eg: tcpdump_xxx_file.pcap.")
	cmd.Flags().StringArrayVar(&cfg.Filter, "filter", cfg.Filter, "The custom filter(UDP,ALL,OnlyTCP).")
	return cmd, nil
}

func run(ctx context.Context, filename string, cfg Config) error {
	if cfg.Verbose {
		cfg.Show[tcpdump.LogDecodeError] = true
	}
	logs := cfg.Logger
	logs.Log(tcpdump.LogDefault, "%s start. file: %s, cfg: %s", color.GreenString("[Tcpdump]"), filename, commons.ToJsonString(cfg))
	options := NewDecodeOptions()
	var source PacketSource
	var err error
	if commons.CheckStdInFromPiped() {
		source = NewConsulSource(os.Stdin, options)
		cfg.Show[tcpdump.LogTCPReassembly] = false
		cfg.Show[tcpdump.LogDecodeError] = true
	} else {
		if source, err = NewFileSource(filename, options, cfg.Loopback); err != nil {
			return err
		}
		cfg.Show[tcpdump.LogTCPReassembly] = true
	}
	newDecoder := tcpdump.NewDefaultDecoder(true, cfg.Logger, map[string]tcpdump.Decoder{
		"HTTP":   tcpdump.NewHTTP1Decoder(),
		"Thrift": tcpdump.NewThriftDecoder(),
	})
	assembler := reassembly.NewAssembler(newDecoder, func(option *reassembly.TCPStreamOption) {
		*option = cfg.TCPStreamOption
	})
	for packet := range source.Packets() {
		if packet == nil {
			continue
		}
		done := func() {
			if wait, isOk := packet.(WaitPacket); isOk {
				wait.Notify()
			}
		}
		success := false

		// tcp
		if tcp, isOK := packet.TransportLayer().(*layers.TCP); isOK {
			if cfg.DisableReassembly {
				header := utils.TCPDumpHeader(tcp, packet.Metadata().Timestamp, nil, utils.NewTCPMetaInfo(packet.NetworkLayer().NetworkFlow(), packet.TransportLayer().TransportFlow(), false), nil)
				logs.Log(tcpdump.LogDefault, header)
				newDecoder().Decode(nil, tcp.LayerPayload())
			} else {
				assembler.AssembleWithContext(packet.NetworkLayer().NetworkFlow(), tcp, &reassembly.Context{CaptureInfo: packet.Metadata().CaptureInfo})
			}
			success = true
		}

		// udp
		if udp, isOk := packet.TransportLayer().(*layers.UDP); isOk && cfg.EnableUDP() {
			if len(udp.LayerPayload()) > 0 {
				hex := codec.NewHexDumpCodec().Encode(udp.LayerPayload())
				logs.Log(tcpdump.LogDefault, string(hex))
			}
			success = true
		}
		if !success && IsFileSource(source) && cfg.EnableALL() {
			dump := packet.Dump()
			logs.Log(tcpdump.LogDefault, dump)
		}
		if !success && IsConsulSource(source) && cfg.EnableALL() {
			payload := packet.TransportLayer().LayerPayload()
			if len(payload) > 0 {
				logs.Log(tcpdump.LogDefault, string(codec.NewHexDumpCodec().Encode(payload)))
			}
		}
		done()
	}
	return nil
}
