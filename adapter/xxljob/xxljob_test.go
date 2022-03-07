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

package xxljob

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xxl-job/xxl-job-executor-go"
)

func TestInitEnabledFalse(t *testing.T) {
	opts := &Options{
		Enabled:      false,
		ServerAddr:   "http://127.0.0.1:19583",
		ExecutorIp:   "10.1.1.123",
		ExecutorPort: "19934",
		LogDir:       "./logs/aaa",
	}
	Init(opts, "hahaha")
	assert.Nil(t, gXxlExecutor)
}

func TestInitServerAddrEmpty(t *testing.T) {
	opts := &Options{
		Enabled:      true,
		ServerAddr:   "",
		ExecutorIp:   "10.1.1.123",
		ExecutorPort: "19934",
		LogDir:       "./logs/aaa",
	}
	err := Init(opts, "hahaha")
	assert.Error(t, err)
}

func TestInitParamsEmpty(t *testing.T) {
	opts := &Options{
		Enabled:      true,
		ServerAddr:   "http://127.0.0.1:19583",
		ExecutorIp:   "",
		ExecutorPort: "",
		LogDir:       "",
	}
	err := Init(opts, "hahaha")
	assert.Nil(t, err)
	//不传配置字段
	opts = &Options{
		Enabled:    true,
		ServerAddr: "http://127.0.0.1:19583",
	}
	err = Init(opts, "hahaha")
	assert.Nil(t, err)
}

func TestInitEnabledTrue(t *testing.T) {
	opts := &Options{
		Enabled:      true,
		ServerAddr:   "http://127.0.0.1:19583",
		ExecutorIp:   "10.1.1.123",
		ExecutorPort: "19934",
		LogDir:       "./logs/aaa",
	}
	Init(opts, "hahaha")
	assert.NotNil(t, gXxlExecutor)
}

func TestNewXxlJobLogger(t *testing.T) {
	opts := &Options{
		Enabled:      true,
		ServerAddr:   "http://127.0.0.1:19583",
		ExecutorIp:   "10.1.1.123",
		ExecutorPort: "19934",
		LogDir:       "./logs/aaa",
	}
	Init(opts, "hahaha")
	_, err := NewXxlJobLogger(0)
	assert.Nil(t, err)
}

func TestRegRask(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("{\"code\":200}"))
	}))
	opts := &Options{
		Enabled:      true,
		ServerAddr:   server.URL,
		ExecutorIp:   "127.0.0.1",
		ExecutorPort: "19934",
		LogDir:       "./logs/aaa",
	}
	Init(opts, "hahaha")
	RegTask("helloworld", func(cxt context.Context, param *xxl.RunReq, logger *XxlJobLogger) string {
		logger.Infof("123")
		logger.Errorf("Errorf test")
		return "ok"
	})

	runReq := &xxl.RunReq{
		JobID:           1,
		LogID:           1,
		ExecutorHandler: "helloworld",
		ExecutorParams:  "123",
	}
	time.Sleep(time.Duration(1) * time.Second)
	bytes, _ := json.Marshal(runReq)
	res, err := http.Post("http://127.0.0.1:19934/run", "application/json", strings.NewReader(string(bytes)))
	if err != nil {
		t.Fatal(err)
	}
	body, _ := ioutil.ReadAll(res.Body)
	var result xxl.LogRes
	json.Unmarshal([]byte(body), &result)
	assert.Equal(t, int64(200), result.Code)
}
