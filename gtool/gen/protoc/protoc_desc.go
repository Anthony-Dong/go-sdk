package protoc

import (
	"fmt"
	"path/filepath"
)

type ProtocGenDesc struct {
	Output string
}

func (p *ProtocGenDesc) init() error {
	if p == nil {
		return nil
	}
	if p.Output == "" {
		return fmt.Errorf(`ProtocGenDesc: the output is nil`)
	}
	if p.Output == "@tmp" {
		if err := setTmpDir(&p.Output); err != nil {
			return err
		}
		p.Output = filepath.Join(p.Output, "proto.desc")
	}
	return nil
}

func (p *ProtocGenDesc) NewDescPlugin() []string {
	if p == nil {
		return []string{}
	}
	return []string{"--descriptor_set_out=" + p.Output}
}
