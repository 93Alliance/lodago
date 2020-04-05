package lodago

// Multimap 允许key值重复
type Multimap struct {
	m    map[interface{}][]interface{}
	size int
}

// NewMultimap 构建multimap
func NewMultimap() *Multimap {
	return &Multimap{
		m:    make(map[interface{}][]interface{}),
		size: 0,
	}
}

// At 取出key对应的values
func (multimap *Multimap) At(key interface{}) ([]interface{}, bool) {
	values, found := multimap.m[key]
	return values, found
}

// Insert 插入一条记录
func (multimap *Multimap) Insert(key interface{}, value interface{}) {
	multimap.m[key] = append(multimap.m[key], value)
	multimap.size++
}

// InsertValues 插入多条记录
func (multimap *Multimap) InsertValues(key interface{}, values []interface{}) {
	for _, v := range values {
		multimap.Insert(key, v)
	}
}

// Remove 移除一条记录
func (multimap *Multimap) Remove(key interface{}, value interface{}) {
	values, found := multimap.m[key]
	if found {
		for idx, v := range values {
			if v == value {
				multimap.m[key] = append(values[:idx], values[idx+1:]...)
				multimap.size--
			}
		}
	}
	if len(multimap.m[key]) == 0 {
		delete(multimap.m, key)
	}
}

// RemoveAll 删除关于key的所有键值对
func (multimap *Multimap) RemoveAll(key interface{}) {
	values, found := multimap.m[key]
	if found {
		multimap.size -= len(values)
		delete(multimap.m, key)
	}
}

// Size 获取map当前键值对的数量
func (multimap *Multimap) Size() int {
	return multimap.size
}

// IsEmpty 判断容器内是否为空
func (multimap *Multimap) IsEmpty() bool {
	return multimap.size == 0
}

// Count 获取key对应的value数量
func (multimap *Multimap) Count(key interface{}) int {
	values, found := multimap.m[key]
	if found {
		return len(values)
	}
	return 0
}
