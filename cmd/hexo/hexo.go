package hexo

import (
	"path/filepath"

	"github.com/anthony-dong/go-sdk/gotool/config"
	"github.com/anthony-dong/go-sdk/gotool/utils"
	"github.com/spf13/cobra"
)

func NewCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{Use: "hexo", Short: "The hexo tool"}
	if err := utils.AddCmd(cmd, newBuildCmd); err != nil {
		return nil, err
	}
	return cmd, nil
}

type hexoConfig struct {
	Dir       string   `json:"dir"`
	Keyword   []string `json:"keyword"`
	Ignore    []string `json:"ignore"`
	TargetDir string   `json:"target_dir"`
}

func newBuildCmd() (*cobra.Command, error) {
	var (
		cfg = &hexoConfig{}
	)
	cmd := &cobra.Command{Use: "build", Short: "Build the markdown project to hexo"}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		commandConfig, err := config.GetCommandConfig(cmd.Context())
		if err != nil {
			return err
		}
		for _, elem := range commandConfig.Hexo.KeyWord {
			cfg.Keyword = append(cfg.Keyword, elem)
		}
		for _, elem := range commandConfig.Hexo.Ignore {
			cfg.Ignore = append(cfg.Ignore, elem)
		}
		if dir, err := filepath.Abs(cfg.Dir); err != nil {
			return err
		} else {
			cfg.Dir = dir
		}
		if dir, err := filepath.Abs(cfg.TargetDir); err != nil {
			return err
		} else {
			cfg.TargetDir = dir
		}
		return buildHexo(cmd.Context(), cfg.Dir, cfg.TargetDir, cfg.Keyword, cfg.Ignore)
	}
	cmd.Flags().StringVarP(&cfg.Dir, "dir", "d", "", "The markdown project dir")
	cmd.Flags().StringVarP(&cfg.Dir, "target_dir", "t", "", "The hexo post dir")
	cmd.Flags().StringArrayVarP(&cfg.Keyword, "keyword", "k", nil, "The keyword need clear, eg: baidu => xxxxx, read from command and load config")
	if err := cmd.MarkFlagRequired("dir"); err != nil {
		return nil, err
	}
	if err := cmd.MarkFlagRequired("target_dir"); err != nil {
		return nil, err
	}
	if err := cmd.MarkFlagRequired("keyword"); err != nil {
		return nil, err
	}
	return cmd, nil
}
