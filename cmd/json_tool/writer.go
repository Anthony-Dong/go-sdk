package json_tool

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/go-sdk/commons"
	"github.com/anthony-dong/go-sdk/commons/logs"
)

func newWriteCli() (*cobra.Command, error) {
	var (
		file   string
		output string
	)
	cmd := &cobra.Command{
		Use:   "writer [--file file] [--output output] ...",
		Short: "Output a file:content json to a dir",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				jsonReader io.Reader
			)
			if output == "" {
				output = commons.MustTmpDir("", commons.GetCmdName())
			}
			if file == "" {
				jsonReader = os.Stdin
			} else {
				file, err := os.Open(file)
				if err != nil {
					return err
				}
				defer file.Close()
				jsonReader = file
			}
			fileMap := make(map[string]string, 0)
			if err := json.NewDecoder(jsonReader).Decode(&fileMap); err != nil {
				return err
			}
			for filename, content := range fileMap {
				filename = filepath.Join(output, filename)
				if err := os.MkdirAll(filepath.Dir(filename), commons.DefaultDirMode); err != nil {
					return err
				}
				logs.Infof("[ReadJsonFile] write file: %s", filename)
				if err := commons.WriteFile(filename, []byte(content)); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", "a json file")
	cmd.Flags().StringVarP(&output, "output", "o", "", "output path")
	return cmd, nil
}
