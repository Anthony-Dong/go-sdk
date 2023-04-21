package codec

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/anthony-dong/go-sdk/commons"

	"google.golang.org/protobuf/proto"
	//"google.golang.org/protobuf/types/pluginpb"
)

func newPBDescCodecCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "pb_desc",
		Short: "decode protobuf-desc file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !commons.CheckStdInFromPiped() {
				return cmd.Help()
			}
			req, err := readReq(os.Stdin)
			if err != nil {
				return err
			}
			fmt.Println(protojson.Format(req))
			return nil
		},
	}
	return cmd, nil
}

func readReq(in io.Reader) (*descriptorpb.FileDescriptorSet, error) {
	//descriptorpb.FileDescriptorSet{}
	req := &descriptorpb.FileDescriptorSet{}
	all, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}
	if err := proto.Unmarshal(all, req); err != nil {
		return nil, err
	}
	return req, nil
}

func writeResp(err error, content map[string]string) {
	ptrStr := func(s string) *string {
		return &s
	}
	resp := &pluginpb.CodeGeneratorResponse{}
	if err != nil {
		resp.Error = ptrStr(err.Error())
	} else {
		for k, v := range content {
			resp.File = append(resp.File, &pluginpb.CodeGeneratorResponse_File{
				Name:    ptrStr(k),
				Content: ptrStr(v),
			})
		}
	}

	marshal, err := proto.Marshal(resp)
	if err != nil {
		panic(err)
	}

	if _, err := os.Stdout.Write(marshal); err != nil {
		panic(err)
	}
}
