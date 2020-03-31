package lodago

import (
	"math/rand"
	"reflect"
	"time"
)

// Find 数组查找函数
func Find(arr interface{}, callback func(ele interface{}, index int) bool) (interface{}, bool) {
	contentValue := reflect.ValueOf(arr)
	arrType := reflect.TypeOf(arr)
	if arrType.Kind() != reflect.Slice {
		panic("expects a slice type")
	}
	for i := 0; i < contentValue.Len(); i++ {
		if content := contentValue.Index(i); callback(content.Interface(), i) {
			return content.Interface(), true
		}
	}
	return nil, false
}

// Filter 切片元素过滤
func Filter(arr interface{}, callback func(ele interface{}, index int) bool) (interface{}, bool) {
	contentType := reflect.TypeOf(arr)
	contentValue := reflect.ValueOf(arr)
	newContent := reflect.MakeSlice(contentType, 0, 0)
	for i := 0; i < contentValue.Len(); i++ {
		if content := contentValue.Index(i); callback(content.Interface(), i) {
			newContent = reflect.Append(newContent, content)
		}
	}
	if newContent.Len() == 0 {
		return newContent.Interface(), false
	} else {
		return newContent.Interface(), true
	}
}

// Shuffle 打乱数组
func Shuffle(arr interface{}) {
	contentType := reflect.TypeOf(arr)
	if contentType.Kind() != reflect.Slice {
		panic("expects a slice type")
	}
	contentValue := reflect.ValueOf(arr)
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	len := contentValue.Len()
	for i := len - 1; i > 0; i-- {
		j := random.Intn(i + 1)
		x, y := contentValue.Index(i).Interface(), contentValue.Index(j).Interface()
		contentValue.Index(i).Set(reflect.ValueOf(y))
		contentValue.Index(j).Set(reflect.ValueOf(x))
	}
}
