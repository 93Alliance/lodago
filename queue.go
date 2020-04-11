package lodago

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"
)

// 本结构参考
// https://github.com/yireyun/go-queue
// https://github.com/golangCasQueue/casQueue

// 缓存
type qcache struct {
	putNo uint64
	getNo uint64
	value interface{}
}

// Queue 无锁队列
type Queue struct {
	capacity  uint64
	capMod    uint64
	putPos    uint64
	getPos    uint64
	cache     []qcache
	sleepTime time.Duration
}

// NewQueue 创建一个队列
func NewQueue(capacity uint64, sleepTime time.Duration) *Queue {
	q := new(Queue)
	// 初始化内部成员变量
	q.capacity = minQuantity(capacity) // TODO: 什么意思？
	q.capMod = q.capacity - 1
	q.putPos = 0
	q.getPos = 0
	q.sleepTime = sleepTime
	q.cache = make([]qcache, q.capacity)
	for i := range q.cache {
		cache := &q.cache[i]
		// 初始化cache内部成员
		cache.putNo = uint64(i)
		cache.getNo = uint64(i)
	}
	cache := &q.cache[0] // 取出一个缓存，设置它的getNo和putNo为本队列的总容量
	cache.getNo = q.capacity
	cache.putNo = q.capacity
	return q
}

// ToString 序列化成字符串
func (q *Queue) ToString() string {
	getPos := atomic.LoadUint64(&q.getPos) // 必须要使用原子操作获取
	putPos := atomic.LoadUint64(&q.putPos) // 必须要使用原子操作获取
	return fmt.Sprintf("Queue{capacity: %v, capMod: %v, putPos: %v, getPos: %v}",
		q.capacity, q.capMod, putPos, getPos)
}

// GetCapacity 获取容量
func (q *Queue) GetCapacity() uint64 {
	return q.capacity
}

// GetQuantity 获取当前队列剩余多少条记录
func (q *Queue) GetQuantity() uint64 {
	quantity := uint64(0)
	getPos := atomic.LoadUint64(&q.getPos)
	putPos := atomic.LoadUint64(&q.putPos)
	if putPos >= getPos { // 如果插入的位置比取出的位置大，那么数量就是插入位置减去取出位置就是剩余的数量了。
		quantity = putPos - getPos
	}
	return quantity
}

// Put 向队列插入数据，返回是否成功，剩余数量。
func (q *Queue) Put(value interface{}) (bool, uint64) {
	var putPos, newPutPos, getPos, posCnt uint64
	var cache *qcache
	getPos = atomic.LoadUint64(&q.getPos)
	putPos = atomic.LoadUint64(&q.putPos)
	// 计算剩余的pos，此pos就可以理解为队列内的一条记录
	if putPos > getPos {
		posCnt = putPos - getPos
	} else {
		posCnt = 0
	}
	if posCnt >= q.capacity { // 如果posCnt大于队列自身容量，说明已经满了。
		time.Sleep(q.sleepTime) // 睡眠一段时间，正常是时间越小越好，但是也要看情况。
		return false, posCnt
	}

	newPutPos = putPos + 1 // 新的putPos
	// 先比较变量的值是否等于给定旧值，等于旧值的情况下才赋予新值，最后返回新值是否设置成功。
	if !atomic.CompareAndSwapUint64(&q.putPos, putPos, newPutPos) {
		runtime.Gosched() // 处理器的时间间隙
		return false, posCnt
	}
	// // newPutPos&q.capMod == newPutPos % q.capacity when q.capacity is 2^n
	cache = &q.cache[newPutPos&q.capMod] // 这个步骤是在做取余操作，相当于分块。
	for {                                // 无限循环
		getNo := atomic.LoadUint64(&cache.getNo)
		putNo := atomic.LoadUint64(&cache.putNo)
		if newPutPos == putNo && getNo == putNo {
			cache.value = value                        // 将值写入队列
			atomic.AddUint64(&cache.putNo, q.capacity) //将缓存内的putNo设置为队列容量
			return true, posCnt + 1
		}
		runtime.Gosched()
	}
}

// Get 从队列中获取记录，返回取出的值，是否成功，剩余数量。
func (q *Queue) Get() (interface{}, bool, uint64) {
	var putPos, getPos, newGetPos, posCnt uint64
	var cache *qcache
	putPos = atomic.LoadUint64(&q.putPos)
	getPos = atomic.LoadUint64(&q.getPos)
	// 计算剩余的pos，此pos就可以理解为队列内的一条记录
	if putPos > getPos {
		posCnt = putPos - getPos
	} else {
		posCnt = 0
	}
	if posCnt < 1 { // 如果剩余小于1也就是等于0，那么失败返回。
		time.Sleep(q.sleepTime)
		return nil, false, posCnt
	}
	newGetPos = getPos + 1
	// 先比较变量的值是否等于给定旧值，等于旧值的情况下才赋予新值，最后返回新值是否设置成功。
	if !atomic.CompareAndSwapUint64(&q.getPos, getPos, newGetPos) {
		runtime.Gosched()
		return nil, false, posCnt
	}

	// putPosNew&q.capMod == putPosNew % q.capacity when q.capacity is 2 ^ n
	cache = &q.cache[newGetPos&q.capMod]
	for {
		getNo := atomic.LoadUint64(&cache.getNo)
		putNo := atomic.LoadUint64(&cache.putNo)
		if newGetPos == getNo && getNo == putNo-q.capacity {
			value := cache.value                       // 取出值
			cache.value = nil                          // 将原有值指向nil
			atomic.AddUint64(&cache.getNo, q.capacity) // 设定缓存的getNo为队列的容量
			return value, true, posCnt - 1
		}
		runtime.Gosched()
	}
}

// Puts 向队列插入多条数据，返回添加的记录数量，剩余数量。
func (q *Queue) Puts(values []interface{}) (int, uint64) {
	var putPos, newPputPos, getPos, posCnt, putCnt uint64
	getPos = atomic.LoadUint64(&q.getPos)
	putPos = atomic.LoadUint64(&q.putPos)
	if putPos > getPos { // 计算剩余的pos，此pos就可以理解为队列内的一条记录
		posCnt = putPos - getPos
	} else {
		posCnt = 0
	}
	// 如果已经满了，就不能添加了。
	if posCnt >= q.capacity {
		time.Sleep(q.sleepTime)
		return 0, posCnt
	}
	if capPuts, size := q.capacity-posCnt, uint64(len(values)); capPuts >= size {
		putCnt = size
	} else {
		putCnt = capPuts
	}
	newPputPos = putPos + putCnt
	if !atomic.CompareAndSwapUint64(&q.putPos, putPos, newPputPos) {
		runtime.Gosched()
		return 0, posCnt
	}
	for posNew, v := putPos+1, uint64(0); v < putCnt; posNew, v = posNew+1, v+1 {
		// putPosNew&q.capMod == putPosNew % q.capacity when q.capacity is 2 ^ n
		var cache = &q.cache[posNew&q.capMod]
		for {
			getNo := atomic.LoadUint64(&cache.getNo)
			putNo := atomic.LoadUint64(&cache.putNo)
			if posNew == putNo && getNo == putNo {
				cache.value = values[v]
				atomic.AddUint64(&cache.putNo, q.capacity)
				break
			} else {
				runtime.Gosched()
			}
		}
	}
	return int(putCnt), posCnt + putCnt
}

// Gets 获取多条记录，返回获取的记录数量，剩余数量。
func (q *Queue) Gets(values []interface{}) (int, uint64) {
	var putPos, getPos, newGetPos, posCnt, getCnt uint64

	putPos = atomic.LoadUint64(&q.putPos)
	getPos = atomic.LoadUint64(&q.getPos)

	if putPos > getPos {
		posCnt = putPos - getPos
	} else {
		posCnt = 0
	}

	if posCnt < 1 {
		time.Sleep(q.sleepTime)
		return 0, posCnt
	}

	if size := uint64(len(values)); posCnt >= size {
		getCnt = size
	} else {
		getCnt = posCnt
	}
	newGetPos = getPos + getCnt

	if !atomic.CompareAndSwapUint64(&q.getPos, getPos, newGetPos) {
		runtime.Gosched()
		return 0, posCnt
	}

	for posNew, v := getPos+1, uint64(0); v < getCnt; posNew, v = posNew+1, v+1 {
		// putPosNew&q.capMod == putPosNew % q.capacity when q.capacity is 2 ^ n
		var cache = &q.cache[posNew&q.capMod]
		for {
			getNo := atomic.LoadUint64(&cache.getNo)
			putNo := atomic.LoadUint64(&cache.putNo)
			if posNew == getNo && getNo == putNo-q.capacity {
				values[v] = cache.value
				cache.value = nil
				getNo = atomic.AddUint64(&cache.getNo, q.capacity)
				break
			} else {
				runtime.Gosched()
			}
		}
	}

	return int(getCnt), posCnt - getCnt
}

// minQuantity 将传入的值转换成2的次方，遵循最小原则，例如：2->2，4->4，7->8,9->16
func minQuantity(v uint64) uint64 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}
