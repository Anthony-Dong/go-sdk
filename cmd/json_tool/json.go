package json_tool

import (
	"github.com/anthony-dong/go-sdk/gotool/utils"
	"github.com/spf13/cobra"
)

func NewCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "json",
		Short: "Json tool",
	}
	if err := utils.AddCmd(cmd, newReaderCmd); err != nil {
		return nil, err
	}
	if err := utils.AddCmd(cmd, newWriteCli); err != nil {
		return nil, err
	}
	return cmd, nil
}
