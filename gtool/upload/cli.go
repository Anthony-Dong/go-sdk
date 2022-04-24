package upload

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/anthony-dong/go-sdk/commons"
	"github.com/anthony-dong/go-sdk/commons/codec"
	"github.com/anthony-dong/go-sdk/commons/logs"
	"github.com/anthony-dong/go-sdk/gtool/config"
)

type uploadCommand struct {
	OssConfigType  string `json:"type"`
	File           string `json:"file"`
	FileNameDecode string `json:"decode"`
}

func NewCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{Use: "upload", Short: `File upload tool`}
	var (
		cfg = &uploadCommand{}
	)
	cmd.Flags().StringVarP(&cfg.File, "file", "f", "", "Set the local path of upload file")
	cmd.Flags().StringVarP(&cfg.FileNameDecode, "decode", "d", "uuid", "Set the upload file name decode method (uuid|url|md5)")
	cmd.Flags().StringVarP(&cfg.OssConfigType, "type", "t", "default", "Set the upload config type")
	if err := cmd.MarkFlagRequired("file"); err != nil {
		return nil, err
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return cfg.Run(cmd.Context())
	}
	return cmd, nil
}

func (c *uploadCommand) validate() error {
	if c.File == "" {
		return fmt.Errorf("flag needs an argument: --file")
	}
	if filename, err := filepath.Abs(c.File); err != nil {
		return err
	} else {
		c.File = filename
	}
	if c.OssConfigType == "" {
		c.OssConfigType = "default"
	}
	logs.Infof("[upload] start config:\n%s", commons.ToPrettyJsonString(c))
	return nil
}
func (c *uploadCommand) Run(ctx context.Context) error {
	if err := c.validate(); err != nil {
		return err
	}
	commandConfig, err := config.GetCommandConfig(ctx)
	if err != nil {
		return err
	}
	if len(commandConfig.Upload.Bucket) == 0 {
		return errors.Errorf("not found bucket config, bucket: %s", c.OssConfigType)
	}
	cfg := commandConfig.Upload.Bucket[c.OssConfigType]
	prefix, suffix := commons.GetFilePrefixAndSuffix(c.File)
	fileInfo := OssUploadFile{
		LocalFile:  c.File,
		SaveDir:    time.Now().Format(commons.FormatTimeV2),
		FilePrefix: c.getFileName(prefix),
		FileSuffix: suffix,
	}
	bucket, err := NewBucket(&cfg)
	if err != nil {
		return err
	}
	if err := fileInfo.PutFile(bucket, &cfg); err != nil {
		return err
	}
	fileUrl := fileInfo.GetOSSUrl(&cfg)
	if logs.IsLevel(logs.LevelInfo) {
		logs.Infof("[upload] end success, url: %s", fileUrl)
	} else {
		fmt.Println(fileUrl)
	}
	return nil
}

func (c *uploadCommand) getFileName(fileName string) string {
	switch c.FileNameDecode {
	case "uuid":
		return commons.GenerateUUID()
	case "url":
		return string(codec.NewUrlCodec().Encode(codec.String2Slice(fileName)))
	case "md5":
		return string(codec.NewMd5Codec().Encode(codec.String2Slice(fileName)))
	}
	return commons.GenerateUUID()
}
