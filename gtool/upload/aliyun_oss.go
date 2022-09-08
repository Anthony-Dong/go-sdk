package upload

import (
	"fmt"
	"os"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/anthony-dong/go-sdk/gtool/config"
)

type OssUploadFile struct {
	LocalFile  string `json:"local_file"` // 本地文件
	SaveDir    string `json:"save_dir"`   // 保存到远程的地址
	FilePrefix string `json:"file_name"`  // 文件名称
	FileSuffix string `json:"file_type"`  // 文件类型名称

	DstFile string `json:"dst_file"` // 目标文件
}

// image/2019-08-29/38564c69-85ba-4415-93d8-xxxxx.jpg.
func (f *OssUploadFile) GetPutPath(config *config.OSSConfig) string {
	if f.DstFile != "" {
		dst := f.DstFile
		dst = strings.TrimLeftFunc(dst, func(r rune) bool {
			return r == '.' || r == '/'
		})
		return fmt.Sprintf("%s/%s", config.PathPrefix, dst)
	}
	return fmt.Sprintf("%s/%s/%s", config.PathPrefix, f.SaveDir, fmt.Sprintf("%s%s", f.FilePrefix, f.FileSuffix))
}

// https://xxxx.oss-accelerate.xxxx.com/image/2020/9-14/xxxxxx.png
func (f *OssUploadFile) GetOSSUrl(config *config.OSSConfig) string {
	path := f.GetPutPath(config)
	return fmt.Sprintf("https://%s/%s", config.UrlEndpoint, path)
}

// NewBucket 创建桶.
func NewBucket(ossConfig *config.OSSConfig) (*oss.Bucket, error) {
	client, err := oss.New(ossConfig.Endpoint, ossConfig.AccessKeyId, ossConfig.AccessKeySecret, func(client *oss.Client) {
		client.Config.Timeout = 5
	})
	if err != nil {
		return nil, err
	}
	// bucket
	bucket, err := client.Bucket(ossConfig.Bucket)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

// PutFile 上传文件.
func (f *OssUploadFile) PutFile(bucket *oss.Bucket, ossConfig *config.OSSConfig) error {
	file, err := os.Open(f.LocalFile)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()
	return bucket.PutObject(f.GetPutPath(ossConfig), file, oss.ObjectStorageClass(oss.StorageStandard), oss.ObjectACL(oss.ACLPublicRead))
}
