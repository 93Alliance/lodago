package lodago

import (
	"reflect"

	json "github.com/json-iterator/go"

	"github.com/mitchellh/mapstructure"
)

// StructToMap 利用反射将结构体转化为map
func StructToMap(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		data[obj1.Field(i).Name] = obj2.Field(i).Interface()
	}
	return data
}

// Map2Struct 将map对象转换成结构体
func Map2Struct(mapValue map[string]interface{}, obj *interface{}) error {
	//将 map 转换为指定的结构体
	if err := mapstructure.Decode(mapValue, obj); err != nil {
		return err
	}
	return nil
}

// Map2JSON map容器转换成json字符串
func Map2JSON(mapValue map[string]interface{}) (string, error) {
	jsonStr, err := json.Marshal(mapValue)
	if err != nil {
		return "", err
	}
	return string(jsonStr), nil
}

// JSON2Map json字符串转换成map容器
func JSON2Map(jsonStr []byte) (map[string]interface{}, error) {
	var mapResult map[string]interface{}
	if err := json.Unmarshal(jsonStr, &mapResult); err != nil {
		return mapResult, err
	}
	return mapResult, nil
}
