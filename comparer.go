package lodago

import "reflect"

// IsEqual 判断两个对象是否相同
func IsEqual(obj1 interface{}, obj2 interface{}) bool {
	return reflect.DeepEqual(obj1, obj2)
}
