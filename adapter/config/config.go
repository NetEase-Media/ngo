// Copyright Ngo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/NetEase-Media/ngo/adapter/log"
	"github.com/spf13/viper"
)

var (
	defaultConfig *Config
	//配置文件目录,可能包含多个配置文件
	gConfigDir string
)

type Config struct {
	*viper.Viper
}

func Init(configName string) (err error) {

	log.Infof("config file path: %s", configName)

	setDefaultConfigDir(configName)

	defaultConfig, err = NewFromFile(configName)
	if err != nil {
		return
	}

	//加载需要 import 的文件
	initConfigImports(defaultConfig, defaultConfig)

	return err
}

// 设置配置文件目录
func setDefaultConfigDir(configName string) (err error) {
	if len(configName) == 0 {
		return errors.New("configName is empty")
	}
	dir, _ := filepath.Split(configName)
	gConfigDir, err = filepath.Abs(dir)

	log.Infof("config base dir: %s", gConfigDir)

	if err != nil {
		return err
	}
	return nil
}

// 初始化configImports
func initConfigImports(rootConfig *Config, currentConfig *Config) {
	imports := currentConfig.GetStringSlice("configImports")
	if len(imports) == 0 {
		return
	}
	for _, file := range imports {
		conf, err := NewFromConfigFile(file)
		if err != nil {
			continue
		}
		if len(conf.GetStringSlice("configImports")) > 0 && filepath.Ext(file) == ".yaml" {
			initConfigImports(rootConfig, conf)
		}
		rootConfig.MergeConfigMap(conf.AllSettings())
	}
}

// NewFromFile 读取配置文件并解析
func NewFromFile(configName string) (config *Config, err error) {
	config = &Config{Viper: viper.New()}
	config.AutomaticEnv()

	configPath, fileName := filepath.Split(configName)
	config.SetConfigName(fileName)
	config.AddConfigPath(configPath)

	// 获取配置类型
	var configType string
	ext := filepath.Ext(fileName)
	if ext != "" {
		ext = ext[1:]
		for _, s := range viper.SupportedExts {
			if s == ext {
				configType = ext
				break
			}
		}
	}
	if configType == "" {
		configType = "yaml"
	}
	config.SetConfigType(configType)
	err = config.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return config, nil
}

// NewFromContent 通过配置文件内容生成配置对象
func NewFromContent(content, configType string) (*Config, error) {
	in := strings.NewReader(content)
	conf := &Config{Viper: viper.New()}
	conf.SetConfigType(configType)
	err := conf.ReadConfig(in)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

// DefaultConfig 返回全局唯一的配置管理对象
func DefaultConfig() *Config {
	return defaultConfig
}

// GetSliceSize 获取大小
func GetSliceSize(key string) int {
	return len(defaultConfig.GetStringSlice(key))
}

// Unmarshal 解析子配置
func Unmarshal(key string, value interface{}) error {
	return defaultConfig.UnmarshalKey(key, value)
}

// 获取配置文件的绝对路径
func GetConfigFilePath(fileName string) string {
	return gConfigDir + string(os.PathSeparator) + fileName
}

//读取配置目录中的配置文件内容
func ReadConfigFile(fileName string) ([]byte, error) {
	return ioutil.ReadFile(GetConfigFilePath(fileName))
}

// 读取viper支持的文件
func NewFromConfigFile(fileName string) (*Config, error) {
	p := GetConfigFilePath(fileName)
	log.Infof("load conf file: %s", p)
	return NewFromFile(GetConfigFilePath(fileName))
}
