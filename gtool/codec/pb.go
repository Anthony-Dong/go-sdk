package codec

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/anthony-dong/go-sdk/commons"
	"github.com/anthony-dong/go-sdk/commons/codec/pb_codec"
	"github.com/anthony-dong/go-sdk/commons/codec/pb_codec/codec"
	"github.com/spf13/cobra"
)

//  echo "CgVoZWxsbxCIBEIDCIgE" | bin/gtool codec base64 --decode | bin/gtool codec pb | jq
func newPBCodecCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "pb",
		Short: "decode protobuf protocol",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !commons.CheckStdInFromPiped() {
				return cmd.Help()
			}
			in, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf(`read std.in find err: %v`, err)
			}
			message, err := pb_codec.DecodeMessage(cmd.Context(), codec.NewBuffer(in))
			if err != nil {
				return fmt.Errorf(`decode pb message find err: %v`, err)
			}
			_, _ = os.Stdout.WriteString(commons.ToJsonString(message))
			return nil
		},
	}
	return cmd, nil
}
