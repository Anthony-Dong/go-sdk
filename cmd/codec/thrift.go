package codec

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/go-sdk/commons"
	"github.com/anthony-dong/go-sdk/commons/codec/thrift_codec"
	"github.com/apache/thrift/lib/go/thrift"
)

//  echo "AAAAEYIhAQRUZXN0HBwWAhUCAAAA" | bin/go-tool codec base64 --decode | bin/go-tool codec thrift | jq
func newThriftCodecCmd() (*cobra.Command, error) {
	messageType := "message"
	cmd := &cobra.Command{
		Use:     "thrift",
		Example: `echo '' | gotool codec thrift`,
		RunE: func(cmd *cobra.Command, args []string) error {
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
			handlerStruct := func(protocol thrift.TProtocol) error {
				data, err := thrift_codec.DecodeMessage(ctx, protocol)
				if err != nil {
					wrapperError(err)
					return err
				}
				_, _ = os.Stdout.WriteString(commons.ToJsonString(data))
				return nil
			}
			switch messageType {
			case "message":
				bufReader := bufio.NewReader(os.Stdin)
				protocol, err := thrift_codec.GetProtocol(bufReader)
				if err != nil {
					return err
				}
				data, err := thrift_codec.DecodeMessage(ctx, thrift_codec.NewTProtocol(bufReader, protocol))
				if err != nil {
					return err
				}
				data.Protocol = protocol
				_, _ = os.Stdout.WriteString(commons.ToJsonString(data))
				return nil
			case "struct":
				if err := handlerStruct(thrift_codec.NewTProtocol(os.Stdin, thrift_codec.UnframedBinary)); err == nil {
					return err
				}
				if err := handlerStruct(thrift_codec.NewTProtocol(os.Stdin, thrift_codec.UnframedCompact)); err == nil {
					return err
				}
			}
			return result
		},
	}
	cmd.Flags().StringVar(&messageType, "type", "", "消息类型, (struct|message)")
	if err := cmd.MarkFlagRequired("type"); err != nil {
		return nil, err
	}
	return cmd, nil
}
