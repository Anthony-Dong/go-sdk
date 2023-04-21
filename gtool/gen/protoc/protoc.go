package protoc

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/iancoleman/orderedmap"
	"github.com/mohae/deepcopy"

	"github.com/anthony-dong/go-sdk/commons"
	"github.com/anthony-dong/go-sdk/commons/logs"
)

const (
	PbSuffix               = ".proto"
	protocCommand          = "protoc"
	protocGenGoCommand     = "protoc-gen-go"
	protocGenGoGrpcCommand = "protoc-gen-go-grpc"
)

var (
	skipFiles = map[string]bool{
		"google/protobuf/any.proto":             true,
		"google/protobuf/api.proto":             true,
		"google/protobuf/compiler/plugin.proto": true,
		"google/protobuf/descriptor.proto":      true,
		"google/protobuf/duration.proto":        true,
		"google/protobuf/empty.proto":           true,
		"google/protobuf/field_mask.proto":      true,
		"google/protobuf/source_context.proto":  true,
		"google/protobuf/struct.proto":          true,
		"google/protobuf/timestamp.proto":       true,
		"google/protobuf/type.proto":            true,
		"google/protobuf/wrappers.proto":        true,
	}
)

type ProtocGen struct {
	RootPath        string
	MainIDL         string
	Include         []string
	NotSkipGooglePb bool // 默认false，表示跳过google pb

	ProtocCommand string // 命令行(生成pb)
	ProtocVersion string
	Go            *ProtocGenGo   `json:"Go,omitempty"`
	Desc          *ProtocGenDesc `json:"Desc,omitempty"`
}

func (p *ProtocGen) init() error {
	var err error
	if p.ProtocCommand == "" {
		if path, err := exec.LookPath(protocCommand); err != nil {
			return err
		} else {
			p.ProtocCommand = path
		}
	}
	p.ProtocVersion, _ = runCmd(p.ProtocCommand, "--version")
	p.MainIDL, err = filepath.Abs(p.MainIDL)
	if err != nil {
		return fmt.Errorf(`filepath.Abs("%s") return err: %v`, p.MainIDL, err)
	}
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	p.RootPath = pwd
	for index, elem := range p.Include {
		abs, err := filepath.Abs(elem)
		if err != nil {
			return fmt.Errorf(`filepath.Abs("%s") return err: %v`, elem, err)
		}
		rel, err := filepath.Rel(p.RootPath, abs)
		if err != nil {
			return fmt.Errorf(`filepath.Rel("%s", "%s") return err: %v`, p.RootPath, abs, err)
		}
		p.Include[index] = rel
	}
	if len(p.Include) == 0 {
		p.Include = []string{"."}
	}
	sortFiles(p.Include)
	return nil
}

func NewProtocGen(main string, ops ...func(gen *ProtocGen)) *ProtocGen {
	config := ProtocGen{
		MainIDL: main,
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
func (p *ProtocGen) Gen(ctx context.Context) error {
	if err := p.init(); err != nil {
		return err
	}
	if err := p.Go.init(); err != nil {
		return err
	}
	if err := p.Desc.init(); err != nil {
		return err
	}
	logs.CtxInfof(ctx, "[Protoc] start, config: %s", commons.ToPrettyJsonString(p))
	files, err := LookupInclude(p.MainIDL, deepcopy.Copy(p.Include).([]string))
	if err != nil {
		return err
	}
	command := make([]string, 0)
	command = append(command, "--experimental_allow_proto3_optional")
	for _, elem := range p.Include {
		command = append(command, fmt.Sprintf("--proto_path=%s", elem))
	}

	command = append(command, p.Go.NewGoPlugin(files, p.RootPath)...)
	command = append(command, p.Desc.NewDescPlugin()...)

	// compile files
	for _, file := range files.Keys() {
		rel, err := filepath.Rel(p.RootPath, file)
		if err != nil {
			return fmt.Errorf(`filepath.Rel("%s", "%s") reture err: %v`, p.RootPath, file, err)
		}
		command = append(command, rel)
	}
	cmd := exec.CommandContext(ctx, p.ProtocCommand, command...)
	fmt.Println(buildCommand(p.ProtocCommand, command))
	return commons.RunWithShell(cmd)
}

var escaperDefault = strings.NewReplacer(`'`, `'\''`)
var notSpaceRegexp = regexp.MustCompile(`^\S*$`)

func buildCommand(exec string, command []string) string {
	builder := strings.Builder{}
	builder.WriteString(exec)
	builder.WriteString(" \\\n")
	for index, elem := range command {
		if notSpaceRegexp.MatchString(elem) {
			builder.WriteString(elem)
		} else {
			builder.WriteByte('\'')
			builder.WriteString(escaperDefault.Replace(elem))
			builder.WriteByte('\'')
		}
		if index == len(command)-1 {
			continue
		}
		builder.WriteString(" \\\n")
	}
	return builder.String()
}

func LookupInclude(main string, include []string) (_ *orderedmap.OrderedMap, err error) {
	defer func() {
		if e := recover(); e != nil {
			switch v := e.(type) {
			case string:
				err = errors.New(v)
			case error:
				err = v
			default:
				err = fmt.Errorf("%#v", e)
			}
		}
	}()
	main, err = filepath.Abs(main)
	if err != nil {
		return nil, err
	}
	for index, elem := range include {
		abs, err := filepath.Abs(elem)
		if err != nil {
			return nil, err
		}
		include[index] = abs
	}
	getRelativeName := func(name string) string {
		for _, elem := range include {
			if strings.HasPrefix(name, elem) {
				rel, _ := filepath.Rel(elem, name)
				return rel
			}
		}
		panic(fmt.Errorf(`not found in include path, file: %s`, name))
	}
	result := orderedmap.New()
	result.Set(main, map[string]bool{
		getRelativeName(main): true,
	})
	if err := handlerFiles(main, include, nil, result); err != nil {
		return nil, err
	}
	return result, nil
}

func handlerFiles(filename string, include []string, walkMap map[string]bool, result *orderedmap.OrderedMap) error {
	if walkMap == nil {
		walkMap = map[string]bool{}
	}
	if walkMap[filename] {
		return nil
	}
	importPaths, err := readImportFile(filename)
	if err != nil {
		return err
	}
	getAbsPath_ := func(name string) (includeName string, abs string) {
		for _, elem := range include {
			absPath := filepath.Join(elem, name)
			if commons.Exist(absPath) {
				return elem, absPath
			}
		}
		panic(fmt.Sprintf("not found file: %s. please check include path.", name))
	}
	walkMap[filename] = true
	for _, elem := range importPaths {
		if skipFiles[elem] {
			continue
		}
		includeName, absImport := getAbsPath_(elem)
		relativeName, err := filepath.Rel(includeName, absImport)
		if err != nil {
			return err
		}
		if value, isExist := result.Get(absImport); isExist {
			value.(map[string]bool)[relativeName] = true
		} else {
			result.Set(absImport, map[string]bool{relativeName: true})
		}
		if err := handlerFiles(absImport, include, deepcopy.Copy(walkMap).(map[string]bool), result); err != nil {
			return err
		}
	}
	return nil
}
