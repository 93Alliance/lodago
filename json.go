package lodago

import (
	"reflect"

	json "github.com/json-iterator/go"
)

// SelectFields 针对性的选择struct字段输出
func SelectFields(v interface{}, fields ...string) ([]byte, error) {
	rt, rv := reflect.TypeOf(v), reflect.ValueOf(v)
	if rt.Kind() != reflect.Ptr || rt.Elem().Kind() != reflect.Struct {
		panic("expects a struct pointer")
	}
	var jsonStr []byte
	var err error
	fs := fieldSet(fields...)
	if len(fs) != 0 { // 如果没有需要排除的字段，那么按照正常序列化执行。
		rt, rv = rt.Elem(), rv.Elem()
		out := make(map[string]interface{}, rt.NumField())
		for i := 0; i < rt.NumField(); i++ {
			field := rt.Field(i)
			jsonKey := field.Tag.Get("json")
			if fs[jsonKey] {
				out[jsonKey] = rv.Field(i).Interface()
			}
		}
		jsonStr, err = json.Marshal(out)
	} else {
		jsonStr, err = json.Marshal(v)
	}
	if err != nil {
		return []byte{}, err
	}
	return jsonStr, nil
}

// DropFields 选择性的抛弃struct的一些字段
func DropFields(v interface{}, fields ...string) ([]byte, error) {
	rt, rv := reflect.TypeOf(v), reflect.ValueOf(v)
	if rt.Kind() != reflect.Ptr || rt.Elem().Kind() != reflect.Struct {
		panic("expects a struct pointer")
	}
	var jsonStr []byte
	var err error
	fs := fieldSet(fields...)
	if len(fs) != 0 { // 如果没有需要排除的字段，那么按照正常序列化执行。
		rt, rv = rt.Elem(), rv.Elem()
		out := make(map[string]interface{}, rt.NumField())
		for i := 0; i < rt.NumField(); i++ {
			field := rt.Field(i)
			jsonKey := field.Tag.Get("json")
			if !fs[jsonKey] {
				out[jsonKey] = rv.Field(i).Interface()
			}
		}
		jsonStr, err = json.Marshal(out)
	} else {
		jsonStr, err = json.Marshal(v)
	}
	if err != nil {
		return []byte{}, err
	}
	return jsonStr, nil
}

// DropMapFields 针对map选择性字段
func DropMapFields(m map[string]interface{}, keys ...string) (map[string]interface{}, error) {
	fs := fieldSet(keys...)
	if len(fs) == 0 {
		return m, nil
	}
	for k := range m {
		if fs[k] {
			delete(m, k) // 删除元素
		}
	}
	return m, nil
}

// 变参扩展
func fieldSet(fields ...string) map[string]bool {
	set := make(map[string]bool, len(fields))
	for _, s := range fields {
		set[s] = true
	}
	return set
}
