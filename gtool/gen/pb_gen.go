package gen

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/anthony-dong/go-sdk/commons"
	"github.com/anthony-dong/go-sdk/commons/logs"
)

const (
	PbSuffix        = ".proto"
	Protoc          = "protoc"
	ProtocGenGo     = "protoc-gen-go"
	ProtocGenGoGrpc = "protoc-gen-go-grpc"
)

var (
	skipFiles = map[string]struct{}{
		"google/protobuf/any.proto":             struct{}{},
		"google/protobuf/api.proto":             struct{}{},
		"google/protobuf/compiler/plugin.proto": struct{}{},
		"google/protobuf/descriptor.proto":      struct{}{},
		"google/protobuf/duration.proto":        struct{}{},
		"google/protobuf/empty.proto":           struct{}{},
		"google/protobuf/field_mask.proto":      struct{}{},
		"google/protobuf/source_context.proto":  struct{}{},
		"google/protobuf/struct.proto":          struct{}{},
		"google/protobuf/timestamp.proto":       struct{}{},
		"google/protobuf/type.proto":            struct{}{},
		"google/protobuf/wrappers.proto":        struct{}{},
	}
)

type ProtocGen struct {
	IDLPath         string
	Include         []string
	BasePackage     string
	OutPutPath      string
	NotSkipGooglePb bool // 默认false，表示跳过google pb

	ProtocCommand          string // 命令行(生成pb)
	ProtocGenGoCommand     string // 生成 全部文件，包含model
	ProtocGenGoGrpcCommand string // 生成 ${service_idl}_grpc.pb.go 文件
}

type ProtocGenResult struct {
	Command  string
	StdOut   bytes.Buffer
	StdError bytes.Buffer
}

func NewProtocGen(dir string, pkg string, ops ...func(gen *ProtocGen)) *ProtocGen {
	config := ProtocGen{
		IDLPath:     dir,
		BasePackage: pkg,
	}
	for _, elem := range ops {
		elem(&config)
	}
	return &config
}

// Gen 首先需要下载
// go-pb: https://github.com/golang/protobuf
// go-pb-doc: https://developers.google.com/protocol-buffers/docs/reference/go-generated#package
// go-grpc: https://github.com/grpc/grpc-go
// go-grpc-doc: https://grpc.io/docs/languages/go/quickstart/
func (p *ProtocGen) Gen(ctx context.Context) (*ProtocGenResult, error) {
	if p.ProtocCommand == "" {
		if path, err := exec.LookPath(Protoc); err != nil {
			return nil, err
		} else {
			p.ProtocCommand = path
		}
	}
	if p.ProtocGenGoCommand == "" {
		if path, err := exec.LookPath(ProtocGenGo); err != nil {
			return nil, err
		} else {
			p.ProtocGenGoCommand = path
		}
	}
	if p.ProtocGenGoGrpcCommand == "" {
		if path, err := exec.LookPath(ProtocGenGoGrpc); err != nil {
			return nil, err
		} else {
			p.ProtocGenGoGrpcCommand = path
		}
	}
	if p.OutPutPath == "" {
		p.OutPutPath = filepath.Join(commons.GetGoPath(), "src")
	}
	output, err := filepath.Abs(p.OutPutPath)
	if err != nil {
		return nil, err
	}
	p.OutPutPath = output
	idlpath, err := filepath.Abs(p.IDLPath)
	if err != nil {
		return nil, err
	}
	p.IDLPath = idlpath

	logs.CtxInfof(ctx, "[Protoc] start, config: %s", commons.ToPrettyJsonString(p))

	files, err := getAllFiles(p.IDLPath, func(fileName string) bool {
		return filepath.Ext(fileName) == PbSuffix
	})
	if err != nil {
		return nil, err
	}
	files = p.filterFiles(files)

	command := make([]string, 0)
	command = append(command, "--experimental_allow_proto3_optional")
	// -I
	command = append(command, fmt.Sprintf("--proto_path=%s", p.IDLPath))
	for _, elem := range p.Include {
		command = append(command, fmt.Sprintf("--proto_path=%s", elem))
	}

	// protobuf plugin
	command = append(command, p.NewPlugin(files, p.ProtocGenGoCommand)...)
	command = append(command, p.NewPlugin(files, p.ProtocGenGoGrpcCommand)...)

	// compile files
	for _, elem := range files {
		rel, _ := filepath.Rel(p.IDLPath, elem)
		command = append(command, rel)
	}
	result := ProtocGenResult{
		Command: buildCommand(p.ProtocCommand, command),
	}

	cmd := exec.CommandContext(ctx, p.ProtocCommand, command...)
	cmd.Stdout = &result.StdOut
	cmd.Stderr = &result.StdError
	if err := cmd.Run(); err != nil {
		return &result, err
	}
	return &result, nil
}

//  NewPlugin file, cmd: /Users/fanhaodong/go/bin/protoc-gen-go
// 解析后 plugin=go
// --plugin=protoc-gen-go=/Users/fanhaodong/go/bin/protoc-gen-go \
// --go_opt=Mapi/api.proto=code.test.org/fanhaodong.516/tool/pb_gen/webcast.room.api/api/api \
// --go_out=/Users/bytedance/go/src
func (p *ProtocGen) NewPlugin(files []string, cmd string) []string {
	tag := strings.TrimPrefix(filepath.Base(cmd), "protoc-gen-")
	command := make([]string, 0)
	// --plugin
	command = append(command, fmt.Sprintf("--plugin=protoc-gen-%s=%s", tag, cmd))
	// --go_opt
	for _, elem := range files {
		goRelativePath, _ := filepath.Rel(p.IDLPath, elem)
		goPackage := p.BasePackage
		if relativePackage := filepath.Dir(goRelativePath); relativePackage != "" && relativePackage != "." {
			goPackage = goPackage + "/" + relativePackage
		}
		command = append(command, fmt.Sprintf("--%s_opt=M%s=%s", tag, goRelativePath, goPackage))
	}
	// --go_out
	command = append(command, fmt.Sprintf("--%s_out=%s", tag, p.OutPutPath))
	return command
}

func buildCommand(exec string, command []string) string {
	builder := strings.Builder{}
	builder.WriteString(exec)
	builder.WriteString(" \\\n")
	for index, elem := range command {
		builder.WriteString(elem)
		if index == len(command)-1 {
			continue
		}
		builder.WriteString(" \\\n")
	}
	return builder.String()
}

// GetAllFiles 从路径dirPth下获取全部的文件.
func getAllFiles(dirPth string, filter func(fileName string) bool) ([]string, error) {
	files := make([]string, 0)
	err := filepath.Walk(dirPth, func(path string, info os.FileInfo, err error) error {
		if info != nil && info.IsDir() {
			return nil
		}
		if filter(path) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (p *ProtocGen) filterFiles(files []string) []string {
	if p.NotSkipGooglePb {
		return files
	}
	result := make([]string, 0, len(files))
	for _, elem := range files {
		rel, _ := filepath.Rel(p.IDLPath, elem)
		if _, isExist := skipFiles[rel]; isExist {
			continue
		}
		result = append(result, elem)
	}
	return result
}
