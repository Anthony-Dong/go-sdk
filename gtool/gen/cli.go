package gen

import (
	"os"

	"github.com/anthony-dong/go-sdk/commons"
	"github.com/anthony-dong/go-sdk/gtool/gen/protoc"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/go-sdk/commons/logs"
	"github.com/anthony-dong/go-sdk/gtool/utils"
)

func NewCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{Use: "gen", Short: `Auto compile thrift、protobuf IDL`}
	if err := utils.AddCmd(cmd, newPBCmd); err != nil {
		return nil, err
	}
	return cmd, nil
}

func newPBCmd() (*cobra.Command, error) {
	var (
		main          string
		include       []string
		plugin        []string
		protocGenGo   protoc.ProtocGenGo
		protocGenDesc protoc.ProtocGenDesc
	)
	cmd := cobra.Command{
		Use:   "protoc [-I include] [--go_pkg package] [--idl idl]",
		Short: `Auto compile protobuf IDL`,
		Long: `Plugin make it easy to compile PB:
Help Doc:
	golang: https://go.dev
	proto3: https://developers.google.com/protocol-buffers/docs/proto3
	protoc-gen-go: https://developers.google.com/protocol-buffers/docs/reference/go-generated
	grpc: https://grpc.io/docs/what-is-grpc/introduction
	protoc-gen-go-grpc: https://grpc.io/docs/languages/go/quickstart
Install:
	golang: https://go.dev/dl
	protoc: https://github.com/protocolbuffers/protobuf/releases
	protoc-gen-go: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	protoc-gen-go-grpc: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			gen := protoc.NewProtocGen(main, func(gen *protoc.ProtocGen) {
				gen.Include = include
				if commons.ContainsString(plugin, "go") {
					gen.Go = &protocGenGo
				}
				if commons.ContainsString(plugin, "desc") {
					gen.Desc = &protocGenDesc
				}
			})
			if err := gen.Gen(ctx); err != nil {
				return err
			}
			logs.CtxInfof(ctx, "exec success")
			return nil
		},
	}
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	cmd.Flags().StringVar(&main, "idl", pwd, "The main idl")
	cmd.Flags().StringArrayVarP(&include, "include", "I", []string{}, "Add an IDL search path for includes, ignore project dir")
	cmd.Flags().StringArrayVar(&plugin, "plugin", []string{"go"}, "The protoc plugin")
	cmd.Flags().StringVar(&protocGenGo.OutPutPath, "go_output", "@tmp", "The output dir. @tmp: sys tmp dir")
	cmd.Flags().StringVarP(&protocGenGo.Package, "go_pkg", "", "", "Define output go package, eg: anthony-dong/proto-example/pb-gen")
	cmd.Flags().BoolVarP(&protocGenGo.SourceRelative, "go_source_relative", "", false, "The output filename is derived from the input filename.")
	cmd.Flags().StringVarP(&protocGenGo.Command, "go_protoc-gen-go", "", "", "The local protoc-gen-go cli")
	cmd.Flags().StringVarP(&protocGenGo.GrpcCommand, "go_protoc-gen-go-grpc", "", "", "The local protoc-gen-go-grpc cli")

	cmd.Flags().StringVar(&protocGenDesc.Output, "desc_output", "@tmp", "The output dir. @tmp: sys tmp dir")
	if err := cmd.MarkFlagRequired("idl"); err != nil {
		return nil, err
	}
	return &cmd, nil
}
