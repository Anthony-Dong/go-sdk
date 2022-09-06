package main

import (
	"context"

	"github.com/anthony-dong/go-sdk/gtool/tcpdump"

	"github.com/anthony-dong/go-sdk/gtool/gen"

	"github.com/anthony-dong/go-sdk/commons/logs"
	"github.com/anthony-dong/go-sdk/gtool/config"

	"github.com/anthony-dong/go-sdk/gtool/codec"
	"github.com/anthony-dong/go-sdk/gtool/hexo"
	"github.com/anthony-dong/go-sdk/gtool/json_tool"
	"github.com/anthony-dong/go-sdk/gtool/upload"
	"github.com/anthony-dong/go-sdk/gtool/utils"
	"github.com/spf13/cobra"
)

const (
	version = "v1.0.4"
)

func main() {
	cmd, err := newRootCmd()
	if err != nil {
		utils.ExitError(err)
	}
	if err := cmd.ExecuteContext(context.Background()); err != nil {
		utils.ExitError(err)
	}
}

func newRootCmd() (*cobra.Command, error) {
	var (
		configFile = ""
		logLevel   = ""
	)
	var rootCmd = &cobra.Command{
		Use:               utils.CliName,
		Version:           version,
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		SilenceErrors:     true, // 屏蔽掉执行错误默认打印日志
		SilenceUsage:      true, // 屏蔽掉执行错误打印help
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if configFile != "" {
				config.SetConfigFile(configFile)
			}
			if logLevel != "" {
				logs.SetLevelString(logLevel)
			}
			return nil
		},
	}
	rootCmd.PersistentFlags().StringVarP(&configFile, "config-file", "", config.GetConfigFile(), "set the config file")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "", "debug", "set the log level in \"fatal|error|warn|info|debug\"")
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
	if err := utils.AddCmd(rootCmd, upload.NewCmd); err != nil {
		return nil, err
	}
	if err := utils.AddCmd(rootCmd, gen.NewCmd); err != nil {
		return nil, err
	}
	if err := utils.AddCmd(rootCmd, tcpdump.NewCmd); err != nil {
		return nil, err
	}
	return rootCmd, nil
}
