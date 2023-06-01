package json_tool

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/anthony-dong/go-sdk/commons"

	"github.com/tidwall/gjson"

	"github.com/anthony-dong/go-sdk/gtool/config"

	"github.com/spf13/cobra"
)

func NewCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "json",
		Short: "The Json tool",
	}
	var (
		pretty bool
		path   string
		reader = os.Stdin
		writer = os.Stdout
	)
	cmd.Flags().StringVar(&path, "path", "", "set specified path")
	cmd.Flags().BoolVarP(&pretty, "pretty", "", false, "set pretty json")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if !commons.CheckStdInFromPiped() {
			return cmd.Help()
		}
		ctx := cmd.Context()
		out, err := readPath(ctx, reader, path)
		if err != nil {
			return err
		}
		if pretty {
			out, err = prettyJson(ctx, out)
			if err != nil {
				return err
			}
		}
		if _, err := writer.Write(out); err != nil {
			return err
		}
		return nil
	}
	return cmd, nil
}

func readPath(ctx context.Context, in io.Reader, path string) ([]byte, error) {
	data, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}
	if path == "" {
		return data, nil
	}
	result := gjson.GetBytes(data, path)
	if !result.Exists() {
		return nil, fmt.Errorf(`not found path: %s`, path)
	}
	return []byte(result.String()), nil
}

func prettyJson(ctx context.Context, src []byte) ([]byte, error) {
	out := bytes.Buffer{}
	cfg, err := config.GetCommandConfig(ctx)
	if err != nil {
		return nil, err
	}
	if cfg.Json.Indent == "" {
		cfg.Json.Indent = "  "
	}
	if err := json.Indent(&out, src, "", cfg.Json.Indent); err != nil {
		return nil, err
	}
	if outBytes := out.Bytes(); len(outBytes) > 0 && outBytes[len(outBytes)-1] == '\n' {
		return outBytes, nil
	}
	out.WriteByte('\n')
	return out.Bytes(), nil
}
