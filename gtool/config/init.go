package config

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/anthony-dong/go-sdk/commons"
	"github.com/anthony-dong/go-sdk/commons/logs"
)

var (
	configFile = ""
	decoder    func(in []byte, out interface{}) (err error)
)

const (
	configFileName = ".gtool.yaml"
)

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	configFile = filepath.Join(pwd, ".config", configFileName)
	if !commons.Exist(configFile) {
		dir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		configFile = filepath.Join(dir, configFileName)
	}
	decoder = func(in []byte, out interface{}) (err error) {
		if err := yaml.Unmarshal(in, out); err == nil {
			return nil
		}
		return json.Unmarshal(in, out)
	}
}

func GetConfigFile() string {
	return configFile
}

func SetConfigFile(file string) {
	if file == "" {
		return
	}
	configFile = file
}

func SetDecoder(d func(in []byte, out interface{}) (err error)) {
	decoder = d
}

func GetCommandConfig(ctx context.Context) (*CommandConfig, error) {
	if !commons.Exist(GetConfigFile()) {
		return &CommandConfig{}, nil
	}
	readFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		logs.CtxErrorf(ctx, "[GetCommandConfig] read config file find err: %v, file: %s", err, configFile)
		return nil, err
	}
	config := new(CommandConfig)
	if err := decoder(readFile, config); err != nil {
		logs.CtxErrorf(ctx, "[GetCommandConfig] decoder config find err: %v, decoder: %v", err, decoder)
		return nil, err
	}
	return config, nil
}
