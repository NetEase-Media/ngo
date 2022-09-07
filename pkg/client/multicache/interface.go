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
