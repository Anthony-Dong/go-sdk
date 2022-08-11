package codec

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/go-sdk/commons"
	"github.com/anthony-dong/go-sdk/commons/codec/thrift_codec"
)

//  echo "AAAAEYIhAQRUZXN0HBwWAhUCAAAA" | bin/gtool codec base64 --decode | bin/gtool codec thrift | jq
func newThriftCodecCmd() (*cobra.Command, error) {
	messageType := "message"
	cmd := &cobra.Command{
		Use:   "thrift",
		Short: "decode thrift protocol",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !commons.CheckStdInFromPiped() {
				return cmd.Help()
			}
			var (
				result       error
				wrapperError = func(err error) {
					if result == nil {
						result = err
						return
					}
					result = fmt.Errorf("%s\n%s", result.Error(), err.Error())
				}
				ctx = cmd.Context()
			)
			handlerStruct := func(r io.Reader, proto thrift_codec.Protocol) error {
				data, err := thrift_codec.DecodeMessage(ctx, thrift_codec.NewTProtocol(r, proto))
				if err != nil {
					wrapperError(err)
					return err
				}
				data.Protocol = proto
				_, _ = os.Stdout.WriteString(commons.ToJsonString(data))
				return nil
			}
			switch messageType {
			case "message":
				bufReader := bufio.NewReader(os.Stdin)
				ctx = thrift_codec.InjectMateInfo(ctx)
				protocol, err := thrift_codec.GetProtocol(ctx, bufReader)
				if err != nil {
					return err
				}
				data, err := thrift_codec.DecodeMessage(ctx, thrift_codec.NewTProtocol(bufReader, protocol))
				if err != nil {
					return err
				}
				data.MetaInfo = thrift_codec.GetMateInfo(ctx)
				data.Protocol = protocol
				_, _ = os.Stdout.WriteString(commons.ToJsonString(data))
				return nil
			case "struct":
				data, _ := ioutil.ReadAll(os.Stdin)
				if err := handlerStruct(bytes.NewReader(data), thrift_codec.FramedBinary); err == nil {
					return nil
				}
				if err := handlerStruct(bytes.NewReader(data), thrift_codec.FramedUnStrictBinary); err == nil {
					return nil
				}
				if err := handlerStruct(bytes.NewReader(data), thrift_codec.UnframedBinary); err == nil {
					return nil
				}
				if err := handlerStruct(bytes.NewReader(data), thrift_codec.UnframedCompact); err == nil {
					return nil
				}
			}
			return result
		},
	}
	cmd.Flags().StringVar(&messageType, "type", "message", "消息类型, (struct|message)")
	return cmd, nil
}
