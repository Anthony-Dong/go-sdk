package config

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/anthony-dong/go-sdk/commons"
	"gopkg.in/yaml.v3"
)

func TestDecode(t *testing.T) {
	fileData, err := ioutil.ReadFile("gotool.yaml")
	if err != nil {
		t.Fatal(err)
	}
	config := new(CommandConfig)
	if err := yaml.Unmarshal(fileData, config); err != nil {
		t.Fatal(err)
	}
	t.Log(commons.ToPrettyJsonString(config))
}
func TestEncode(t *testing.T) {
	config := CommandConfig{
		Upload: UploadConfig{
			Bucket: map[string]OSSConfig{
				"default": {
					AccessKeyId:     "xxxx",
					AccessKeySecret: "xxx",
					Endpoint:        "oss-accelerate.xxxxx.com",
					UrlEndpoint:     "xxxx.oss-accelerate.xxxx.com",
					Bucket:          "tyut",
					PathPrefix:      "file",
				},
				"image": {
					AccessKeyId:     "xxxx",
					AccessKeySecret: "xxxx",
					Endpoint:        "oss-accelerate.xxxxx.com",
					UrlEndpoint:     "xxxx.oss-accelerate.xxxx.com",
					Bucket:          "tyut",
					PathPrefix:      "image",
				},
			},
		},
		Hexo: HexoConfig{
			Ignore:  []string{"/bin", ".git"},
			KeyWord: []string{"xxx", "xxx"},
		},
	}
	marshal, err := yaml.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(marshal))
	if err := ioutil.WriteFile(filepath.Join(commons.GetGoProjectDir(), "command/config/gotool.yaml"), marshal, commons.DefaultFileMode); err != nil {
		t.Fatal(err)
	}
}
