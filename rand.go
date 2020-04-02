package lodago

import (
	"math/rand"
	"time"

	uuid "github.com/satori/go.uuid"
)

// UUID 生成唯一的uuid
func UUID() string {
	uuid := uuid.NewV4()
	return uuid.String()
}

// 数字 + 大小写字母
const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	alphabetIdxBits = 6
	alphabetIdxMask = 1<<alphabetIdxBits - 1
	alphabetIdxMax  = 63 / alphabetIdxBits
)

var src = rand.NewSource(time.Now().UnixNano())

// RandString 随机字符串（数字 + 大小写字母）
func RandString(n ...int) string {
	num := 64
	if len(n) > 0 {
		num = n[0]
	}
	b := make([]byte, num)
	for i, cache, remain := num-1, src.Int63(), alphabetIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), alphabetIdxMax
		}
		if idx := int(cache & alphabetIdxMask); idx < len(alphabet) {
			b[i] = alphabet[idx]
			i--
		}
		cache >>= alphabetIdxBits
		remain--
	}
	return Bytes2String(b)
}
