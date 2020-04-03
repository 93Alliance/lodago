package lodago

import "reflect"

// IsFunc 是否为函数
func IsFunc(fun interface{}) bool {
	rt := reflect.TypeOf(fun)
	if rt.Kind() != reflect.Func {
		return false
	}
	return true
}

// After 在执行几次后才真正执行
func After(n int, fun interface{}) interface{} {
	rt := reflect.TypeOf(fun)
	if rt.Kind() != reflect.Func {
		panic("expects a function")
	}
	rv := reflect.ValueOf(fun)
	wrapper := reflect.MakeFunc(rt, func(args []reflect.Value) []reflect.Value {
		n--
		if n < 1 {
			return rv.Call(args)
		}
		results := make([]reflect.Value, rt.NumOut())
		for i := 0; i < rt.NumOut(); i++ {
			results[i] = reflect.Zero(rt.Out(i))
		}
		return results
	})
	return wrapper.Interface()
}

// Before 在执行几次之后不再执行
func Before(n int, fun interface{}) interface{} {
	rt := reflect.TypeOf(fun)
	if rt.Kind() != reflect.Func {
		panic("expects a function")
	}
	rv := reflect.ValueOf(fun)
	wrapper := reflect.MakeFunc(rt, func(args []reflect.Value) []reflect.Value {
		n--
		if n > 0 {
			return rv.Call(args)
		}
		// 当n小于0将不需要执行，需要构造零值返回。
		results := make([]reflect.Value, rt.NumOut())
		for i := 0; i < rt.NumOut(); i++ {
			results[i] = reflect.Zero(rt.Out(i))
		}
		return results
	})
	return wrapper.Interface()
}
