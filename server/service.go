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
	"errors"

	"github.com/NetEase-Media/ngo/adapter/log"
)

var serviceOptions ServiceOptions // TODO: 这个可能需要放到server内

// ServiceOptions 整合了服务的全局配置
type ServiceOptions struct {
	AppName     string
	ClusterName string
	Instance    string
}

// Check 检查配置合法性
func (o *ServiceOptions) Check() error {
	log.Infof("check ServiceOptions: %+v", o)

	if o.AppName == "" || o.ClusterName == "" {
		return errors.New("lack of service config")
	}
	o.Instance = "default"
	return nil
}
