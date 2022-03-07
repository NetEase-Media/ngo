//+build !race

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

package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/NetEase-Media/ngo/adapter/log"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// testGetConfigPath 获取测试配置的路径
func testGetConfigPath() string {
	_, filename, _, _ := runtime.Caller(1)
	fmt.Println(runtime.Caller(1))
	return path.Join(path.Dir(filename), "../configs/config_sample.yaml")
}

func TestInfo(t *testing.T) {
	opt := NewDefaultOptions()
	s := newServer(opt)
	s.AddRoute(GET, "/route1", testRoute1)
	go s.Start()
	s.stopServer(context.Background())
}

func testCheck(c *gin.Context) {
	c.String(http.StatusOK, "test check")
}

func testRoute1(c *gin.Context) {
	c.String(http.StatusOK, "testRoute1")
}

func TestGoAttach(t *testing.T) {
	opt := NewDefaultOptions()
	opt.Port = 9090
	s := newServer(opt)
	go s.Start()
	stopped := false
	s.GoAttach(func() {
		select {
		case <-s.StoppingNotify():
			stopped = true
			return
		}
	})
	s.stopServer(context.Background())
	assert.True(t, stopped)
}

func TestServer(t *testing.T) {
	dir, err := os.Getwd()
	path := dir + string(os.PathSeparator) + "app.yaml"
	content :=
		`
service:
  appName: ngo
  clusterName: ngo-online
httpServer:
  port: 8080
  mode: debug
log:
  - name: default
    level: info
    packageLevel:
      - gorm.io/gorm: error
    path: ./log
    errorPath: ./error
    fileName: ngo
    writableStack: false
    format: txt
    noFile: true
    filePathPattern:
    maxAge: 240h
    rotationTime: 24h
    rotationSize: 2048
pprof:
  switch: true
  port: 8899
`
	ioutil.WriteFile(path, []byte(content), 0666)
	assert.NoError(t, err)

	configPath = path
	s := Init()

	os.Remove(path)

	s.PreStart = func() error {
		log.Info("do pre-start...")
		return nil
	}

	s.PreStop = func(ctx context.Context) error {
		log.Info("do pre-stop...")
		return nil
	}
	go s.Start()
	time.Sleep(time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err = s.Stop(ctx)
	assert.NoError(t, err)
}
