package codec

import (
	"fmt"
	"io"
	"os"

	"github.com/anthony-dong/go-sdk/commons"

	"github.com/anthony-dong/go-sdk/commons/codec"
	"github.com/anthony-dong/go-sdk/gtool/utils"
	"github.com/spf13/cobra"
)

var (
	isDecode bool
)

func NewCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "codec",
		Short: "Encode and decode data",
	}
	cmd.PersistentFlags().BoolVar(&isDecode, "decode", false, "decode content data")
	if err := utils.AddCmd(cmd, newThriftCodecCmd); err != nil {
		return nil, err
	}
	cmd.AddCommand(newCodecCmd("gizp", codec.NewGzipCodec()))
	cmd.AddCommand(newCodecCmd("base64", codec.NewCodec(codec.NewBase64Codec())))
	cmd.AddCommand(newCodecCmd("br", codec.NewBrCodec()))
	cmd.AddCommand(newCodecCmd("snappy", codec.NewCodec(codec.NewSnappyCodec())))
	cmd.AddCommand(newCodecCmd("md5", codec.NewCodec(codec.NewMd5Codec())))
	cmd.AddCommand(newCodecCmd("url", codec.NewCodec(codec.NewUrlCodec())))
	cmd.AddCommand(newCodecCmd("hex", codec.NewCodec(codec.NewHexCodec())))
	return cmd, nil
}

var (
	reader io.Reader = os.Stdin
	writer io.Writer = os.Stdout
)

func newCodecCmd(name string, codec codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("%s codec", name),
		//Flags: codecFlags,
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
}
