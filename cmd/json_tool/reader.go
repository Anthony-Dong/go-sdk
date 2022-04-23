package json_tool

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"

	"github.com/tidwall/gjson"
)

func newReaderCmd() (*cobra.Command, error) {
	path := ``
	cmd := &cobra.Command{
		Use:   "reader [--path path]",
		Short: "Get searches json for the specified path",
		Example: `Exec: echo '{"k1":{"k2":"v2"}}' | bam tool json --path k1.k2
Output: v2

Exit status:
0: if OK
1: if not found specified path
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				return err
			}
			result := gjson.GetBytes(data, path)
			if !result.Exists() {
				return fmt.Errorf(`not found path: %s`, path)
			}
			if _, err := os.Stdout.Write([]byte(result.String())); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&path, "path", "", "json 路径")
	return cmd, nil
}
