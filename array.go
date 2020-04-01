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
	}
	return newContent.Interface(), true
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

// Concat 拼接切片
func Concat(arr1 interface{}, arr2 ...interface{}) interface{} {
	resultType := reflect.TypeOf(arr1)
	if resultType.Kind() != reflect.Slice {
		panic("expects a slice type")
	}
	len := 0
	arr1V := reflect.ValueOf(arr1)
	arr1Len := arr1V.Len()
	len += arr1Len
	for _, s := range arr2 { // 获取切片的长度
		rv := reflect.ValueOf(s)
		len += rv.Len()
	}
	result := reflect.MakeSlice(resultType, len, len) // 提前开辟足够的容量
	index := 0
	for i := 0; i < arr1Len; i++ { // 拷贝被拼接的切片元素
		result.Index(index).Set(arr1V.Index(i))
		index++
	}
	for _, s := range arr2 {
		rv := reflect.ValueOf(s)
		for i := 0; i < rv.Len(); i++ {
			result.Index(index).Set(rv.Index(i))
			index++
		}
	}
	return result.Interface()
}

// Fill 填充切片
func Fill(arr interface{}, value interface{}, options ...int) {
	rt, rv := reflect.TypeOf(arr), reflect.ValueOf(arr)
	valueType := reflect.TypeOf(value)
	if rt.Kind() != reflect.Slice {
		panic("expects a slice")
	}
	if valueType.Kind() != rt.Elem().Kind() {
		panic("expects fill value is " + rt.Elem().Name())
	}
	optionLen := len(options)
	start, end, arrLen := 0, 0, rv.Len()
	if optionLen >= 2 {
		start, end = options[0], options[1]
		if end > arrLen {
			end = arrLen
		}
	} else if optionLen == 1 {
		start, end = options[0], arrLen
	} else {
		end = arrLen
	}
	if start > arrLen-1 {
		start = arrLen - 1
	}
	for i := start; i < end; i++ {
		rv.Index(i).Set(reflect.ValueOf(value))
	}
}

// Remove 移除索引位置的元素
func Remove(arr interface{}, index int, order ...bool) {
	rt, rv := reflect.TypeOf(arr), reflect.ValueOf(arr)
	if rt.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Slice {
		panic("expects a slice pointer")
	}
	rvPtr := rv.Elem()
	arrLen := rvPtr.Len()
	if index >= arrLen {
		return
	}
	if len(order) > 0 && !order[0] {
		// 方式2 不保序
		rvPtr.Index(index).Set(rvPtr.Index(arrLen - 1))
		rvPtr.Index(arrLen - 1).Set(reflect.Zero(rt.Elem().Elem()))
		rvPtr.Set(rvPtr.Slice(0, arrLen-1))
	} else {
		// 方式1 保序
		reflect.Copy(rvPtr.Slice(index, arrLen), rvPtr.Slice(index+1, arrLen))
		rvPtr.Index(arrLen - 1).Set(reflect.Zero(rt.Elem().Elem()))
		rvPtr.Set(rvPtr.Slice(0, arrLen-1))
	}
}
