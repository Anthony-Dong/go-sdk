package codec

import (
	"fmt"
	"io"
	"os"

	"github.com/anthony-dong/go-sdk/commons"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/go-sdk/commons/codec"
	"github.com/anthony-dong/go-sdk/gtool/utils"
)

func NewCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "codec",
		Short: "The Encode and Decode data tool",
	}
	if err := utils.AddCmd(cmd, newThriftCodecCmd); err != nil {
		return nil, err
	}
	if err := utils.AddCmd(cmd, newPBCodecCmd); err != nil {
		return nil, err
	}
	cmd.AddCommand(newCodecCmd("gizp", codec.NewGzipCodec()))
	cmd.AddCommand(newCodecCmd("base64", codec.NewCodec(codec.NewBase64Codec())))
	cmd.AddCommand(newCodecCmd("br", codec.NewBrCodec()))
	cmd.AddCommand(newCodecCmd("deflate", codec.NewDeflateCodec()))
	cmd.AddCommand(newCodecCmd("snappy", codec.NewCodec(codec.NewSnappyCodec())))
	cmd.AddCommand(newCodecCmd("md5", codec.NewCodec(codec.NewMd5Codec())))
	cmd.AddCommand(newCodecCmd("url", codec.NewCodec(codec.NewUrlCodec())))
	cmd.AddCommand(newCodecCmd("hex", codec.NewCodec(codec.NewHexCodec())))
	cmd.AddCommand(newCodecCmd("hexdump", codec.NewCodec(codec.NewHexDumpCodec())))
	return cmd, nil
}

var (
	reader   io.Reader = os.Stdin
	writer   io.Writer = os.Stdout
	isDecode bool
)

func newCodecCmd(name string, codec codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("%s codec", name),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if !commons.CheckStdInFromPiped() {
				return cmd.Help()
			}
			if isDecode {
				if err := codec.Decode(reader, writer); err != nil {
					return err
				}
				return nil
			}
			if err := codec.Encode(reader, writer); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.PersistentFlags().BoolVar(&isDecode, "decode", false, "decode content data")
	return cmd
}
