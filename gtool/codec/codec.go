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
	cmd.AddCommand(newCodecCmd("gzip", codec.NewGzipCodec()))
	cmd.AddCommand(newCodecCmd("base64", codec.NewCodec(codec.NewBase64Codec())))
	cmd.AddCommand(newCodecCmd("br", codec.NewBrCodec()))
	cmd.AddCommand(newCodecCmd("deflate", codec.NewDeflateCodec()))
	cmd.AddCommand(newCodecCmd("snappy", codec.NewCodec(codec.NewSnappyCodec())))
	cmd.AddCommand(newCodecCmd("md5", codec.NewCodec(codec.NewMd5Codec())))
	cmd.AddCommand(newCodecCmd("url", codec.NewCodec(codec.NewUrlCodec())))
	cmd.AddCommand(newCodecCmd("hex", codec.NewCodec(codec.NewHexCodec())))
	cmd.AddCommand(newCodecCmd("hexdump", codec.NewCodec(codec.NewHexDumpCodec())))
	cmd.AddCommand(newCodecCmd("double-quote", codec.NewCodec(bytesCodec{
		BytesEncoder: codec.NewStringQuoteCodec(),
		BytesDecoder: nil,
	})))
	cmd.AddCommand(newCodecCmd("single-quote", codec.NewCodec(bytesCodec{
		BytesEncoder: func() codec.BytesEncoder {
			r := codec.NewStringQuoteCodec()
			r.QuoteType = codec.SingleQuoteClike
			return r
		}(),
		BytesDecoder: BytesDecoderFunc(func(src []byte) (dst []byte, err error) {
			if err := cmd.Help(); err != nil {
				return nil, err
			}
			return nil, fmt.Errorf(`not support decode type`)
		}),
	})))
	cmd.AddCommand(newCodecCmd("pb-desc", codec.NewCodec(bytesCodec{
		BytesEncoder: BytesEncodeFunc(func(src []byte) (dst []byte) {
			return src
		}),
		BytesDecoder: codec.NewProtoDesc(),
	})))
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
			defer func() {
				if err == nil {
					_, _ = writer.Write([]byte{'\n'})
				}
			}()
			if isDecode {
				return codec.Decode(reader, writer)
			}
			return codec.Encode(reader, writer)
		},
	}
	cmd.PersistentFlags().BoolVar(&isDecode, "decode", false, "decode content data")
	return cmd
}

type bytesCodec struct {
	codec.BytesEncoder
	codec.BytesDecoder
}

type BytesDecoderFunc func(src []byte) (dst []byte, err error)
type BytesEncodeFunc func(src []byte) (dst []byte)

func (b BytesDecoderFunc) Decode(src []byte) (dst []byte, err error) {
	return b(src)
}
func (b BytesEncodeFunc) Encode(src []byte) (dst []byte) {
	return b(src)
}
