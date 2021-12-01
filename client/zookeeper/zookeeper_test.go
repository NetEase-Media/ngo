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
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	ADDR = "kafka"
	NAME = "zktest"
)

func TestGetZkClient(t *testing.T) {
	c := GetZkClient("m")
	t.Log(c)
}

func TestZookeeperProxy_CreateNode_flag0(t *testing.T) {
	var opts []Options
	o := Options{
		Name: NAME,
		//Addr: []string{testZkClientAddr},
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 10,
	}
	opts = append(opts, o)
	zookeeperClients, err := NewClientFromOption(opts)
	assert.Equal(t, nil, err, "初始化失败")
	c := zookeeperClients[NAME]
	// 创建节点
	dd := c.Exist("/cc")
	if dd == true {
		ddd := c.Delete("/cc")
		assert.Equal(t, true, ddd, "删除节点失败")
	}
	b := c.CreateNode("/cc", 0, "aa")
	assert.Equal(t, true, b, "创建失败")
	if b == true {
		c.Delete("/cc")
	}
}

func TestZookeeperProxy_CreateNode_flag1(t *testing.T) {
	var opts []Options
	o := Options{
		Name: NAME,
		//Addr: []string{testZkClientAddr},
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 5,
	}
	opts = append(opts, o)
	zookeeperClients, err := NewClientFromOption(opts)
	assert.Equal(t, nil, err, "初始化失败")
	c := zookeeperClients[NAME]

	// 创建节点
	dd := c.Exist("/cc")
	if dd == true {
		ddd := c.Delete("/cc")
		assert.Equal(t, true, ddd, "删除节点失败")
	}
	b := c.CreateNode("/cc", 1, "aa")
	assert.Equal(t, true, b, "创建失败")
	/*dd = c.Exist("/cc")
	assert.Equal(t, false, dd, "节点未自动失效")*/
}

func TestZookeeperProxy_CreateNode_flag2(t *testing.T) {
	var opts []Options
	o := Options{
		Name: NAME,
		//Addr: []string{testZkClientAddr},
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 15,
	}
	opts = append(opts, o)
	zookeeperClients, err := NewClientFromOption(opts)
	assert.Equal(t, nil, err, "初始化失败")
	c := zookeeperClients[NAME]
	// 创建节点(带序号)
	bb := c.CreateNode("/testzk_ss1s_", 2, "flag2")
	assert.Equal(t, true, bb, "创建失败")
	str, _, _ := c.Conn.Children("/")
	for i := range str {
		if strings.Contains(str[i], "testzk_ss1s_") {
			c.Delete("/" + str[i])
		}
	}
}

func TestZookeeperProxy_SetData_pathErr(t *testing.T) {
	var opts []Options
	o := Options{
		Name: NAME,
		//Addr: []string{testZkClientAddr},
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 15,
	}
	opts = append(opts, o)
	zookeeperClients, err := NewClientFromOption(opts)
	assert.Equal(t, nil, err, "初始化失败")
	c := zookeeperClients[NAME]
	// 创建节点(带序号)
	bb := c.CreateNode("/testzk_ss1s_", 2, "flag2")
	assert.Equal(t, true, bb, "创建失败")
	str, _, _ := c.Conn.Children("/")
	for i := range str {
		if strings.Contains(str[i], "testzk_ss1s_") {
			q := c.SetData(str[i], "testzk_ss1s_setData")
			assert.Equal(t, false, q, "设置节点值失败")
			c.Delete("/" + str[i])
			break
		}
	}
}

func TestZookeeperProxy_SetData(t *testing.T) {
	var opts []Options
	o := Options{
		Name: NAME,
		//Addr: []string{testZkClientAddr},
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 15,
	}
	opts = append(opts, o)
	zookeeperClients, err := NewClientFromOption(opts)
	assert.Equal(t, nil, err, "初始化失败")
	c := zookeeperClients[NAME]
	// 创建节点(带序号)
	bb := c.CreateNode("/testzk_ss1s_", 2, "flag2")
	assert.Equal(t, true, bb, "创建失败")
	str, _, _ := c.Conn.Children("/")
	for i := range str {
		if strings.Contains(str[i], "testzk_ss1s_") {
			q := c.SetData("/"+str[i], "testzk_ss1s_setData")
			assert.Equal(t, true, q, "设置节点值失败")
			c.Delete("/" + str[i])
			break
		}
	}
}

func TestZookeeperProxy_Delete(t *testing.T) {
	var opts []Options
	o := Options{
		Name: NAME,
		//Addr: []string{testZkClientAddr},
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 15,
	}
	opts = append(opts, o)
	zookeeperClients, err := NewClientFromOption(opts)
	assert.Equal(t, nil, err, "初始化失败")
	c := zookeeperClients[NAME]
	dd := c.Exist("/testzk_ss0s")
	if dd == true {
		c.Delete("/testzk_ss0s")
	}
	// 创建节点
	bb := c.CreateNode("/testzk_ss0s", 0, "flag0")
	assert.Equal(t, true, bb, "创建失败")
	res := c.Delete("/testzk_ss0s")
	assert.True(t, res)
}

func TestZookeeperProxy_Delete_PathErr(t *testing.T) {
	var opts []Options
	o := Options{
		Name: NAME,
		//Addr: []string{testZkClientAddr},
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 15,
	}
	opts = append(opts, o)
	zookeeperClients, err := NewClientFromOption(opts)
	assert.Equal(t, nil, err, "初始化失败")
	c := zookeeperClients[NAME]
	// 创建节点
	dd := c.Exist("/testzk_ss0s")
	if dd == true {
		c.Delete("/testzk_ss0s")
	}
	bb := c.CreateNode("/testzk_ss0s", 0, "flag0")
	assert.Equal(t, true, bb, "创建失败")
	res := c.Delete("testzk_ss0s")
	assert.False(t, res)
	res = c.Delete("/testzk_ss0s")
	assert.Equal(t, true, res, "删除失败")
}

func TestAll(t *testing.T) {
	var opts []Options
	o := Options{
		Name: "mm",
		//Addr: []string{testZkClientAddr},
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 10,
	}
	opts = append(opts, o)
	zookeeperClients, err := NewClientFromOption(opts)
	//err := Init(opts)
	assert.Equal(t, nil, err, "初始化失败")
	c := zookeeperClients["mm"]

	// 创建节点
	dd := c.Exist("/cc")
	if dd == true {
		ddd := c.Delete("/cc")
		assert.Equal(t, true, ddd, "创建节点失败")
	}
	b := c.CreateNode("/cc", 0, "aa")
	assert.Equal(t, true, b, "创建失败")
	if b == true {
		c.Delete("/cc")
	}
	b = c.CreateNode("/cc", 1, "aa1")
	assert.Equal(t, true, b, "创建失败")

	// 创建节点(带序号)
	bb := c.CreateNode("/ss1s_", 2, "aa")
	assert.Equal(t, true, bb, "创建失败")

	// 检查节点是否存在
	d := c.Exist("/cc")
	assert.Equal(t, true, d, "该节点不存在")

	// 设置节点值
	dddd := c.Exist("/bb")
	if dddd == false {
		sign := c.CreateNode("/bb", 0, "ss")
		assert.Equal(t, true, sign, "")
	}
	q := c.SetData("/bb", "ahs")
	assert.Equal(t, true, q, "设置节点值失败")

	// 删除节点
	ds := c.Exist("/qq")
	if ds == false {
		sign1 := c.CreateNode("/qq", 0, "ss")
		assert.Equal(t, true, sign1, "")
	}
	p := c.Delete("/qq")
	assert.Equal(t, true, p, "删除节点失败")

}

// 监视子节点
func TestZkChildNodeWatcher(t *testing.T) {

	var opts []Options

	o := Options{
		Name: "cc",
		Addr: []string{ADDR},
		//Addr:           []string{"localhost:2181"},
		SessionTimeout: time.Second * 5,
	}
	opts = append(opts, o)
	zookeeperClients, err := NewClientFromOption(opts)
	assert.Equal(t, nil, err, "初始化失败")
	c := zookeeperClients["cc"]

	path := "/bb"

	dd := c.Exist(path)
	if dd == false {
		bb := c.CreateNode(path, 0, "aa")
		assert.Equal(t, true, bb, "创建失败")
	}

	go c.ZkPathChildrenWatcher(path, func(respchan <-chan *NodeChildResponse) {
		select {
		case r := <-respchan:
			//oldchil := r.OldChildren
			newchil := r.NewChildren
			response := r.Response
			//assert.Equal(t, []string{}, oldchil, "获取旧子节点失败")
			assert.Equal(t, []string{"qq"}, newchil, "获取新子节点失败")
			assert.Equal(t, "NodeChildren changed", response, "响应失败")
		}
	})

	pathc := "/bb/qq"
	ddd := c.Exist(pathc)
	if ddd == false {
		// 创建节点
		bbbb := c.CreateNode(pathc, 0, "aa")
		assert.Equal(t, true, bbbb, "创建失败")
	} else {
		cc := c.Delete(pathc)
		assert.Equal(t, true, cc, "删除失败")
		// 创建节点
		bbbb := c.CreateNode(pathc, 0, "aa")
		assert.Equal(t, true, bbbb, "创建失败")
	}

	time.Sleep(10 * time.Second)
}

// 监视指定节点
func TestZkNodeWatcher(t *testing.T) {

	var opts []Options

	o := Options{
		Name: "mm",
		//Addr: []string{testZkClientAddr},
		Addr:           []string{ADDR},
		SessionTimeout: time.Second * 5,
	}
	opts = append(opts, o)
	zookeeperClients, err := NewClientFromOption(opts)
	assert.Equal(t, nil, err, "初始化失败")
	c := zookeeperClients["mm"]

	path := "/qqq"
	dbb := c.Exist(path)
	if dbb == false {
		sign := c.CreateNode(path, 0, "ss")
		assert.Equal(t, true, sign, "")
	}
	go c.ZkNodeWatcher(path, func(respChan <-chan *NodeResponse) {
		select {
		case cc := <-respChan:
			//a := cc.OldData
			b := cc.NewData
			c := cc.Response
			//assert.Equal(t, "juan111134", a, "响应失败")
			assert.Equal(t, "juan111135", b, "响应失败")
			assert.Equal(t, "This Node "+path+" Node Data Changed", c, "响应失败")
		}
	})

	path_ := "/qqq/qqqq"
	ds := c.Exist(path_)
	if ds == false {
		sign1 := c.CreateNode(path_, 0, "ss")
		assert.Equal(t, true, sign1, "")
	}
	go c.ZkNodeWatcher(path_, func(respChan <-chan *NodeResponse) {
		select {
		case cc := <-respChan:
			//a := cc.OldData
			b := cc.NewData
			c := cc.Response
			//assert.Equal(t, "juan1", a, "响应失败")
			assert.Equal(t, "", b, "响应失败")
			assert.Equal(t, "This Node "+path_+" Node was deleted", c, "删除失败")
		}
	})

	// 改变节点值
	q := c.SetData(path, "juan111135")
	assert.Equal(t, true, q, "设置节点值失败")

	// 删除节点
	p := c.Delete(path_)
	assert.Equal(t, true, p, "删除节点值失败")
	time.Sleep(7 * time.Second)

}
