package protoc

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/iancoleman/orderedmap"
)

type ProtocGenGo struct {
	Package    string // go package
	OutPutPath string // go output

	Command        string // protoc-gen-go command
	Version        string // protoc-gen-go version
	GrpcCommand    string // protoc-gen-go-grpc command
	GrpcVersion    string // protoc-gen-go-grpc version
	SourceRelative bool   // 默认false.
}

func (p *ProtocGenGo) init() error {
	if p == nil {
		return nil
	}
	var err error
	if p.OutPutPath == "" {
		return fmt.Errorf(`not found go output path`)
	}
	if p.Package == "" {
		return fmt.Errorf(`not found go package`)
	}
	if p.Command == "" {
		if path, err := exec.LookPath(protocGenGoCommand); err != nil {
			return err
		} else {
			p.Command = path
		}
	}
	p.Version, _ = runCmd(p.Command, "--version")
	if p.GrpcCommand == "" {
		if path, err := exec.LookPath(protocGenGoGrpcCommand); err != nil {
			return err
		} else {
			p.GrpcCommand = path
		}
	}
	p.GrpcVersion, _ = runCmd(p.GrpcCommand, "--version")
	if p.OutPutPath == "@tmp" {
		return setTmpDir(&p.OutPutPath)
	} else {
		if p.OutPutPath, err = filepath.Abs(p.OutPutPath); err != nil {
			return fmt.Errorf(`abs "%s" find err: %v`, p.OutPutPath, err)
		}
		if err := os.MkdirAll(p.OutPutPath, 0755); err != nil {
			return fmt.Errorf(`os.MkdirAll("%s", 0755) return err: %v`, p.OutPutPath, err)
		}
	}
	return nil
}

//  NewGoPlugin file, cmd: /Users/fanhaodong/go/bin/protoc-gen-go
// 解析后 plugin=go
// --plugin=protoc-gen-go=/Users/fanhaodong/go/bin/protoc-gen-go \
// --go_opt=Mapi/api.proto=code.test.org/fanhaodong.516/tool/pb_gen/webcast.room.api/api/api \
// --go_opt=paths=source_relative \
// --go_out=$(go env GOPATH)/src
func (p *ProtocGenGo) NewGoPlugin(files *orderedmap.OrderedMap, rootPath string) []string {
	if p == nil {
		return []string{}
	}
	return append(p.newGoPlugin(files, p.Command, rootPath), p.newGoPlugin(files, p.GrpcCommand, rootPath)...)
}
func (p *ProtocGenGo) newGoPlugin(files *orderedmap.OrderedMap, cmd string, rootPath string) []string {
	tag := strings.TrimPrefix(filepath.Base(cmd), "protoc-gen-")
	command := make([]string, 0)
	// --plugin
	command = append(command, fmt.Sprintf("--plugin=protoc-gen-%s=%s", tag, cmd))
	// --go_opt
	for _, file := range files.Keys() {
		goRelativePath, err := filepath.Rel(rootPath, file)
		if err != nil {
			panic(fmt.Sprintf(`filepath.Rel("%s", "%s") find err: %v`, rootPath, file, err))
		}
		goPackage := p.Package
		if relativePackage := filepath.Dir(goRelativePath); relativePackage != "" && relativePackage != "." {
			goPackage = goPackage + "/" + relativePackage
		}
		includePathList, _ := files.Get(file)
		for path := range includePathList.(map[string]bool) {
			command = append(command, fmt.Sprintf("--%s_opt=M%s=%s", tag, path, goPackage))
		}
	}
	if p.SourceRelative { // source relative 表示 idl路径=code-gen路径, 默认是根据option go_package = ""进行确认的.
		command = append(command, fmt.Sprintf("--%s_opt=%s", tag, "paths=source_relative"))
	}
	// --go_out
	command = append(command, fmt.Sprintf("--%s_out=%s", tag, p.OutPutPath))
	return command
}
