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

package zookeeper

import (
	"github.com/go-zookeeper/zk"
)

// 获取zookeeper客户端
func GetZkClient(name string) *ZookeeperProxy {
	return zookeeperClients[name]
}

// 创建节点
func (z *ZookeeperProxy) CreateNode(path string, flags int32, s string) bool {
	data := []byte(s)
	acls := zk.WorldACL(zk.PermAll) // 获取访问控制权限
	// flags有4种取值:
	//   0:永久,除非手动删除
	//   1:短暂,session断开则改节点也被删除
	//   2:会自动在节点后面添加序号
	//   3:即,短暂且自动添加序号
	_, err := z.Conn.Create(path, data, flags, acls)
	if err != nil {
		return false
	} else {
		return true
	}
}

// 判断节点是否存在
func (z *ZookeeperProxy) Exist(path string) bool {
	sign, _, _ := z.Conn.Exists(path)
	return sign
}

// 设置节点值
func (z *ZookeeperProxy) SetData(path string, s string) bool {
	data := []byte(s)
	_, stat, err := z.Conn.Get(path)
	if err != nil {
		return false
	}
	_, err = z.Conn.Set(path, data, stat.Version)
	if err != nil {
		return false
	} else {
		return true
	}
}

// 删除节点
func (z *ZookeeperProxy) Delete(path string) bool {
	_, stat, err := z.Conn.Get(path)
	if err != nil {
		return false
	}
	err_ := z.Conn.Delete(path, stat.Version)
	if err_ != nil {
		return false
	}
	return true
}

// 用来反馈监听子节点的结果
type NodeChildResponse struct {
	OldChildren []string // 子节点改变前的所有子节点
	NewChildren []string // 子节点改变后的所有子节点
	Response    string   // 监听响应
	Error       error    // 连接错误标志
}

// 循环监听子节点                        TODO: 由于第三方包没有实现相应的功能的支持，因此未完全实现对子节点的动态监听。 需要后续进行改进
func (z *ZookeeperProxy) ZkPathChildrenWatcher(path string, listener func(respChan <-chan *NodeChildResponse)) {
	for {
		respChan := make(chan *NodeChildResponse, 1)
		oldChildren, _, child_ch, err := z.Conn.ChildrenW(path)
		if err != nil {
			respChan <- &NodeChildResponse{[]string{}, []string{}, "Can Not Watch The Children", err}
		}
		select {
		case e := <-child_ch:
			if e.Type == zk.EventNodeChildrenChanged {
				newChildren, _, err := z.Conn.Children(path)
				respChan <- &NodeChildResponse{oldChildren, newChildren, "NodeChildren changed", err}
			}

			listener(respChan)
		}
	}
}

// 用来反馈监听当前节点的结果
type NodeResponse struct {
	OldData  string // 节点状态改变前的数据
	NewData  string // 节点状态改变后的数据
	Response string // 监听响应
	Error    error  // 连接错误标志
}

// 循环监听当前节点
func (z *ZookeeperProxy) ZkNodeWatcher(path string, listener func(respChan <-chan *NodeResponse)) {

	for {
		respChan := make(chan *NodeResponse, 1)
		oldData, _, events, err := z.Conn.GetW(path)
		if err != nil {
			respChan <- &NodeResponse{"", "", "Can Not Watch This Node " + path, err}
		}
		select {
		case e := <-events:
			if e.Type == zk.EventNodeDataChanged {
				newData, _, err := z.Conn.Get(path)
				if err != nil {
					respChan <- &NodeResponse{string(oldData), "", "Can Not Get New Data", err}
				} else {
					respChan <- &NodeResponse{string(oldData), string(newData), "This Node " + path + " Node Data Changed", nil}
				}
			}

			if e.Type == zk.EventNodeDeleted {
				respChan <- &NodeResponse{string(oldData), "", "This Node " + path + " Node was deleted", nil}
			}
			listener(respChan)
		}
	}
}
