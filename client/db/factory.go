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

package db

import (
	"github.com/NetEase-Media/ngo/adapter/log"
)

func Init(opts []*Options) error {
	if len(opts) == 0 {
		// 如果没有db配置则跳过
		log.Info("empty db config, so skip init")
		return nil
	}

	for _, opt := range opts {
		c, err := NewClient(opt)
		if err != nil {
			return err
		}
		dbClients[opt.Name] = c
	}
	return nil
}

var dbClients = make(map[string]*Client, 4)

func GetMysqlClient(name string) *Client {
	return dbClients[name]
}

func GetClient(name string) *Client {
	return dbClients[name]
}

func GetAllClients() []*Client {
	clients := make([]*Client, 0, len(dbClients))
	for _, v := range dbClients {
		clients = append(clients, v)
	}
	return clients
}

func CloseAllClients() {
	clients := GetAllClients()
	for _, v := range clients {
		v.Close()
	}
}