package config

type CommandConfig struct {
	Upload UploadConfig `yaml:"Upload" json:"upload"`
	Hexo   HexoConfig   `yaml:"Hexo" json:"hexo"`
	Json   JsonConfig   `json:"json" yaml:"json"`
}

type UploadConfig struct {
	Bucket map[string]OSSConfig `yaml:"Bucket" json:"bucket"`
}

type JsonConfig struct {
	Indent string `json:"indent" yaml:"indent"` // 默认是两个空格
}
type OSSConfig struct {
	AccessKeyId     string `yaml:"AccessKeyId" json:"access_key_id"`
	AccessKeySecret string `yaml:"AccessKeySecret" json:"access_key_secret"`
	Endpoint        string `yaml:"Endpoint" json:"endpoint"`
	UrlEndpoint     string `yaml:"UrlEndpoint" json:"url_endpoint"`
	Bucket          string `yaml:"Bucket" json:"bucket"`
	PathPrefix      string `yaml:"PathPrefix" json:"path_prefix"`
}

type HexoConfig struct {
	Ignore  []string `yaml:"Ignore,omitempty" json:"ignore"`
	KeyWord []string `yaml:"KeyWord,omitempty" json:"key_word"`
}
