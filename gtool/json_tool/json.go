package json_tool

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/anthony-dong/go-sdk/commons"

	"github.com/tidwall/gjson"

	"github.com/anthony-dong/go-sdk/gtool/config"

	"github.com/anthony-dong/go-sdk/gtool/utils"
	"github.com/spf13/cobra"
)

func NewCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "json",
		Short: "The Json tool",
		Example: fmt.Sprintf(`  Exec: echo '{"k1":{"k2":"v2"}}' | %s json --path k1 --pretty
  Output: {
             "k2": "v2"
          }
  Help: https://github.com/tidwall/gjson`, utils.CliName),
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
	if err := utils.AddCmd(cmd, newWriteCli); err != nil {
		return nil, err
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
	if outBytes := out.Bytes(); outBytes[len(outBytes)-1] == '\n' {
		return outBytes, nil
	}
	out.WriteByte('\n')
	return out.Bytes(), nil
}

func trimIllegalLine(data []byte) ([]byte, error) {
	scanner := bufio.NewScanner(bytes.NewBuffer(data))
	begin := false
	out := bytes.NewBuffer(make([]byte, 0, len(data)))
	writeLine := func(data []byte) error {
		out.Write(data)
		out.WriteByte('\n')
		return nil
	}
	for scanner.Scan() {
		bytesLine := scanner.Bytes()
		if begin {
			if err := writeLine(bytesLine); err != nil {
				return nil, err
			}
		}
		if line := strings.TrimSpace(string(bytesLine)); len(line) > 0 && (line[0] == '{' || line[0] == '[') {
			begin = true
			if err := writeLine(bytesLine); err != nil {
				return nil, err
			}
		}
	}
	result := out.Bytes()
	if len(result) == 0 {
		return result, nil
	}
	return result[:len(result)-1], nil
}
