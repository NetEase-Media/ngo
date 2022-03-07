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

package job

import (
	"context"
	"fmt"
	"time"

	"github.com/NetEase-Media/ngo/adapter/log"
	"github.com/NetEase-Media/ngo/adapter/util"
	"github.com/NetEase-Media/ngo/client/httplib"
)

const (
	maxRetries = 3
)

type Args struct {
	SharedNum int
}

type Callback func(*Args) (string, error)

type Options struct {
	CenterUrl string
	Namespace string
}

func (o *Options) check() {
	if o.CenterUrl == "" {
		log.Fatalf("empty center url")
	}
	if o.Namespace == "" {
		o.Namespace = "default"
	}
}

type job struct {
	opt      *Options
	f        Callback
	hostname string
}

func (t *job) run() {
	// 准备环境，获取分片信息
	res := t.setUp()

	var retString string
	var err error

	// 结束处理
	startTime := time.Now()
	defer t.tearDown(startTime, &retString, &err)

	args := &Args{
		SharedNum: res.Data.SharedNum,
	}
	retString, err = t.f(args)
}

func (t *job) setUp() *GetShareNumResponse {
	var res *GetShareNumResponse
	var err error
	for i := 0; i < maxRetries; i++ {
		res, err = t.getArgs()
		if err == nil {
			break
		}
	}

	util.CheckError(err)
	return res
}

func (t *job) tearDown(startTime time.Time, retString *string, retErr *error) {
	req := &ReportRequest{
		PodName:   t.hostname,
		StartTime: startTime.UnixNano() / int64(time.Millisecond),
		EndTime:   time.Now().UnixNano() / int64(time.Millisecond),
		Result:    *retString,
	}
	if r := recover(); r != nil {
		req.Exception = fmt.Sprintf("%v", r)
	} else {
		if (*retErr) != nil {
			req.Exception = (*retErr).Error()
		}
	}

	log.WithFields(
		"podName", req.PodName,
		"startTime", req.StartTime,
		"endTime", req.EndTime,
		"result", req.Result,
		"exception", req.Exception,
	).Info("report job")

	var err error
	for i := 0; i < maxRetries; i++ {
		err = t.report(req)
		if err == nil {
			break
		}
	}
	if err != nil {
		log.Errorf("report failed: %s", err.Error())
	}
}

type GetShareNumResponse struct {
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Data    GetShareNumData `json:"data"`
}

type GetShareNumData struct {
	SharedNum int `json:"sharedNum"`
}

func (t *job) getArgs() (*GetShareNumResponse, error) {
	url := fmt.Sprintf("%s/api/v1/%s/podshared/getSharedNum", t.opt.CenterUrl, t.opt.Namespace)
	var res GetShareNumResponse
	_, err := httplib.Get(url).AddQuery("podName", t.hostname).BindJson(&res).Do(context.Background())
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ReportRequest 是结果上报结构
type ReportRequest struct {
	PodName   string `json:"podName"`
	Exception string `json:"exception"`
	StartTime int64  `json:"startTime"`
	EndTime   int64  `json:"endTime"`
	Result    string `json:"result"`
}

type ReportResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (t *job) report(req *ReportRequest) error {
	url := fmt.Sprintf("%s/api/v1/%s/podshared/report", t.opt.CenterUrl, t.opt.Namespace)
	var res ReportResponse
	_, err := httplib.Post(url).SetJson(req).BindJson(&res).Do(context.Background())
	return err
}
