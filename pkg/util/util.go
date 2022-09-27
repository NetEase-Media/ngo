package util

import (
	"reflect"
)

// MaxInt64 返回数组中最大的
func MaxInt64(i int64, nums ...int64) int64 {
	max := i
	for _, num := range nums {
		if max < num {
			max = num
		}
	}
	return max
}

// MinInt64 返回数组中最小的
func MinInt64(i int64, nums ...int64) int64 {
	min := i
	for _, num := range nums {
		if min > num {
			min = num
		}
	}
	return min
}

func TypeName(data interface{}) string {
	if data == nil {
		return ""
	}
	t := reflect.TypeOf(data)
	for t.Kind() == reflect.Ptr { // 解引用嵌套指针
		t = t.Elem()
	}

	return t.String()
}
