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

package multicache

// 文件主要定义各级handler处理数据的方式
// 目前只是提供字符串数据的传输和输入，复杂类型不做处理

// Handler 处理方式
type Handler interface {
	Priority() int // 定义Handler的处理优先级0最高，数字越大优先级越低
	Set(key, value string) (bool, error)
	SetWithTimeout(key, value string, ttl int) (bool, error)
	Get(key string) (string, error)
	Evict(key string) (bool, error)
	Clear() (bool, error)
}
