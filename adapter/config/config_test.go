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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var yamlExample = []byte(`
name: 
  steve
hobbies:
- skateboarding
- snowboarding
- go
clothing:
  jacket: leather
  trousers: adsad
  size: large
age: 35
eyes : brown
beard: true
slice:
  - aaa: aaa
    bbb: bbb
  - aaa: aaa
    bbb: bbb
`)

type clothing struct {
	Jacket   string `yaml:jacket,omitempty`
	Trousers string `yaml:trousers,omitempty`
	Size     string `yaml:size,omitempty`
}

func createFile(fileName string, content []byte) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	tempDirName, err := ioutil.TempDir(dir, "hello")
	if err != nil {
		return "", err
	}

	configName := tempDirName + fileName
	ioutil.WriteFile(configName, yamlExample, 0666)
	return configName, nil
}

func TestInitConfigFile(t *testing.T) {
	path, err := createFile("/hello.yaml", yamlExample)
	if err != nil {
		t.Fatalf("createFile error:%v", err)
		return
	}
	defer os.RemoveAll(filepath.Dir(path))

	// Init()
	Init(path)
	conf := DefaultConfig()
	assert.Equal(t, "steve", conf.GetString("name"))
	clothing := &clothing{}
	err = Unmarshal("clothing", clothing)
	if err != nil {
		t.Fatalf("error:%v", err)
	}
	assert.Equal(t, "leather", clothing.Jacket)

	s := make([]string, 0, 8)
	fmt.Println(len(s))
	fmt.Println(len(defaultConfig.GetStringSlice("slice")))
	err = Unmarshal("slice", &s)
	fmt.Println(len(s))
}

func TestInitConfigDir(t *testing.T) {
	gomonkey.ApplyFunc(NewFromFile, func(configName string) (*Config, error) {
		c := &Config{Viper: viper.New()}
		c.Set("app", "hell")
		if strings.HasSuffix(configName, "a.yaml") {
			return c, nil
		}
		return c, nil
	})
	assert.Empty(t, defaultConfig.GetString("app"))
	Init("a.yaml")
	assert.Equal(t, "hell", defaultConfig.GetString("app"))
	//path, err := createFile("/app.yaml", yamlExample)
	//if err != nil {
	//	t.Fatalf("createFile error:%v", err)
	//	return
	//}
	//defer os.RemoveAll(filepath.Dir(path))
	//
	//dir := filepath.Dir(filepath.Dir(path))
	//t.Logf("path:%s, \ndir:%s\n", path, dir)
	//err = Init("", dir)
	//if err != nil {
	//	t.Fatalf("init error,%v", err)
	//}
	//conf := DefaultConfig()
	//assert.Equal(t, "steve", conf.GetString("name"))
	//clothing := &clothing{}
	//err = Unmarshal("clothing", clothing)
	//if err != nil {
	//	t.Fatalf("error:%v", err)
	//}
	//assert.Equal(t, "leather", clothing.Jacket)
}

func TestInitConfigImports(t *testing.T) {
	rootConfig := &Config{Viper: viper.New()}
	rootConfig.Viper.Set("configImports", []string{"a.yaml", "b.yaml"})
	gomonkey.ApplyFunc(NewFromFile, func(configName string) (*Config, error) {
		c := &Config{Viper: viper.New()}
		if strings.HasSuffix(configName, "a.yaml") {
			c.Set("a", "aa")
		} else if strings.HasSuffix(configName, "b.yaml") {
			c.Set("b", "bb")
		}
		return c, nil
	})
	assert.Empty(t, rootConfig.GetString("a"))
	assert.Empty(t, rootConfig.GetString("b"))
	initConfigImports(rootConfig, rootConfig)
	assert.Equal(t, "aa", rootConfig.GetString("a"))
	assert.Equal(t, "bb", rootConfig.GetString("b"))
}

func TestNewFromContent(t *testing.T) {
	config, err := NewFromContent(string(yamlExample), "yaml")
	if err != nil {
		t.Fatalf("error:%v", err)
	}
	assert.Equal(t, "steve", config.GetString("name"))
}

func TestGetConfigFilePath(t *testing.T) {
	gConfigDir = "config"
	fileName := "haha.yaml"
	sep := string(os.PathSeparator)
	expected := fmt.Sprintf("%s%s%s", gConfigDir, sep, fileName)
	assert.Equal(t, expected, GetConfigFilePath(fileName))
}

func TestFileExt(t *testing.T) {
	aext := filepath.Ext("aa.yaml")
	bext := filepath.Ext("bb.yml")
	assert.Equal(t, ".yaml", aext)
	match, err := regexp.MatchString("^\\.ya?ml$", aext)
	if err == nil {
		assert.True(t, match, "ok")
	} else {
		t.Fatalf("fail:%v+", err)
	}
	match, err = regexp.MatchString("^\\.ya?ml$", bext)
	if err == nil {
		assert.True(t, match, "ok")
	} else {
		t.Fatalf("fail:%v+", err)
	}
}

//func TestNewFromFileException(t *testing.T) {
//	//配置文件不存在
//	_, err := NewFromFile("config/elephant-recommend-go/xff.txt")
//	assert.Error(t, err)
//	/*//配置文件后缀
//	path, err := createFile("/hello.txt", yamlExample)
//	if err != nil {
//		t.Fatalf("createFile error:%v", err)
//		return
//	}
//	_,err = NewFromFile(path)*/
//}
