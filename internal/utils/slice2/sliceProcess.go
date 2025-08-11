package slice2

import (
	"fmt"
)

type sliceError struct {
	msg string
}

func (e *sliceError) Error() string {
	return e.msg
}

func Errorf(format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return &sliceError{msg}
}
func RemoveDuplicateStrings(originals []string) []string {
	temp := map[string]struct{}{}
	result := make([]string, 0, len(originals))
	for _, item := range originals {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
func RemoveDuplicateElement2[T fmt.Stringer](originals []T) []T {
	// 切片去重
	// 通过switch选择切片类型，可自定义，原理是将切片值格式化成string，作为map的key，如果不存在则添加到map，判断是否存在map来去重。

	temp := map[string]struct{}{}
	result := make([]T, 0, len(originals))
	for _, item := range originals {
		if _, ok := temp[item.String()]; !ok {
			temp[item.String()] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func RemoveDuplicateElement(originals interface{}) (interface{}, error) {
	// 切片去重
	// 通过switch选择切片类型，可自定义，原理是将切片值格式化成string，作为map的key，如果不存在则添加到map，判断是否存在map来去重。
	temp := map[string]struct{}{}
	switch slice := originals.(type) {
	case [][]string:
		result := make([][]string, 0)
		for _, item := range slice {
			key := fmt.Sprint(item)
			if _, ok := temp[key]; !ok {
				temp[key] = struct{}{}
				result = append(result, item)
			}
		}
		return result, nil
	case []string:
		result := make([]string, 0, len(originals.([]string)))
		for _, item := range slice {
			key := fmt.Sprint(item)
			if _, ok := temp[key]; !ok {
				temp[key] = struct{}{}
				result = append(result, item)
			}
		}
		return result, nil
	case []int64:
		result := make([]int64, 0, len(originals.([]int64)))
		for _, item := range slice {
			key := fmt.Sprint(item)
			if _, ok := temp[key]; !ok {
				temp[key] = struct{}{}
				result = append(result, item)
			}
		}
		return result, nil
	default:
		err := Errorf("Unknown type: %T", slice)
		return nil, err
	}
}
