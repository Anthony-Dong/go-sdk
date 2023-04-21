package tcpdump

import (
	"github.com/anthony-dong/go-sdk/commons"
	"github.com/anthony-dong/go-sdk/commons/tcpdump"
	"github.com/anthony-dong/go-sdk/gtool/tcpdump/reassembly"
)

type Config struct {
	reassembly.TCPStreamOption
	Loopback          bool                    `json:"loopback"`           //  IsLoopback: 表示抓包的网卡是 本地回环网卡，一般就是 lo0 或者 lo
	DisableReassembly bool                    `json:"disable_reassembly"` //  UseReassembly: 表示不进行包重组
	Filter            []string                `json:"filter"`             // Filter
	Verbose           bool                    `json:"verbose"`            //
	Show              map[tcpdump.LogTag]bool `json:"-"`
}

//func (c Config) EnableTCP() bool {
//	return commons.ContainsString(c.Filter, "TCP") || c.EnableALL()
//}

func (c Config) EnableUDP() bool {
	return commons.ContainsString(c.Filter, "UDP") || c.EnableALL()
}

func (c Config) EnableALL() bool {
	return commons.ContainsString(c.Filter, "ALL")
}

func (c Config) OnlyTCP() bool {
	return commons.ContainsString(c.Filter, "OnlyTCP")
}

func NewDefaultConfig() Config {
	config := Config{
		TCPStreamOption:   reassembly.NewDefaultTCPStreamOption(),
		Loopback:          false,
		DisableReassembly: false,
		Show: map[tcpdump.LogTag]bool{
			tcpdump.LogDefault:    true,
			tcpdump.LogDecodeDump: true,
		},
		Filter: []string{"OnlyTCP"},
	}
	config.Logger = tcpdump.NewDefaultLogger(config.Show)
	return config
}
