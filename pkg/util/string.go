package util

import (
	"encoding/json"
	"strconv"
	"strings"
)

// Symbol 分隔符
type Symbol string

// 标记
const (
	Point             Symbol = "."
	Comma             Symbol = ","
	Colon             Symbol = ":"
	Empty             Symbol = ""
	Blank             Symbol = " "
	Underline         Symbol = "_"
	DoubleVerticalBar Symbol = "||"
	VerticalBar       Symbol = "|"
	QuestionMark      Symbol = "?"
	Percent           Symbol = "%"
	Enter             Symbol = "\n"
)

// Join 通过分割符号将数组转换成字符串
func Join(list []string, symbol Symbol) string {
	if len(list) == 0 {
		return string(Empty)
	}
	var l []string
	for _, v := range list {
		// if v == "" { // 过滤掉为空的数据
		// 	continue  // 先去掉兼容java工程
		// }
		l = append(l, v)
	}
	r := strings.Join(l, string(symbol))
	return strings.Trim(r, string(Blank))
}

// Split 将字符串根据分割符号转换成数组
func Split(s string, symbol Symbol) []string {
	var r []string = make([]string, 0)

	if s == "" {
		return r
	}
	l := strings.Split(s, string(symbol))
	if len(l) == 0 {
		return r
	}
	for _, v := range l {
		v = strings.Trim(v, " ")
		if v == "" {
			continue
		}
		r = append(r, v)
	}
	return r
}

// SplitNoRepeat 无重复的数组
func SplitNoRepeat(s string, symbol Symbol) []string {
	var r []string = make([]string, 0)

	if s == "" {
		return r
	}

	l := strings.Split(s, string(symbol))
	if len(l) == 0 {
		return r
	}

	var checker map[string]int = make(map[string]int, len(l))
	for _, v := range l {
		v = strings.Trim(v, " ")
		if v == "" {
			continue
		}
		_, ok := checker[v]
		if ok {
			continue
		}
		checker[v] = 0 // 设置检查位
		r = append(r, v)
	}
	return r
}

// Strval 获取变量的字符串值
// 浮点型 3.0将会转换成字符串3, "3"
// 非数值或字符类型的变量将会被转换成JSON格式字符串
func Strval(value interface{}) string {
	var key string
	if value == nil {
		return key
	}
	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}
	return key
}
