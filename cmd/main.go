package main

import (
	"context"

	"github.com/anthony-dong/go-sdk/gotool/hexo"

	"github.com/anthony-dong/go-sdk/gotool/codec"
	"github.com/anthony-dong/go-sdk/gotool/json_tool"
	"github.com/anthony-dong/go-sdk/gotool/utils"
	"github.com/spf13/cobra"
)

const (
	version = "v1.0.0"
)

func newRootCmd() (*cobra.Command, error) {
	var rootCmd = &cobra.Command{
		Use:               utils.CliName,
		Version:           version,
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		SilenceErrors:     true, // 屏蔽掉执行错误默认打印日志
		SilenceUsage:      true, // 屏蔽掉执行错误打印help
	}
	rootCmd.SetHelpTemplate(utils.HelpTmpl)
	rootCmd.SetUsageTemplate(utils.UsageTmpl)
	if err := utils.AddCmd(rootCmd, codec.NewCmd); err != nil {
		return nil, err
	}
	if err := utils.AddCmd(rootCmd, json_tool.NewCmd); err != nil {
		return nil, err
	}
	if err := utils.AddCmd(rootCmd, hexo.NewCmd); err != nil {
		return nil, err
	}
	return rootCmd, nil
}

func main() {
	cmd, err := newRootCmd()
	if err != nil {
		utils.ExitError(err)
	}
	if err := cmd.ExecuteContext(context.Background()); err != nil {
		utils.ExitError(err)
	}
}
