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
const numBytes = "0123456789"
const lowerBytes = "abcdefghijklmnopqrstuvwxyz"
const upperBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	alphabetIdxBits = 6
	alphabetIdxMask = 1<<alphabetIdxBits - 1
	alphabetIdxMax  = 63 / alphabetIdxBits
	numIdxBits      = 4
	numIdxMask      = 1<<numIdxBits - 1
	numIdxMax       = 63 / numIdxBits
	lowerIdxBits    = 5
	lowerIdxMask    = 1<<lowerIdxBits - 1
	lowerIdxMax     = 63 / lowerIdxBits
	upperIdxBits    = 5
	upperIdxMask    = 1<<upperIdxBits - 1
	upperIdxMax     = 63 / upperIdxBits
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

// RandStrWithNum 随机生成数字字符串
func RandStrWithNum(n ...int) string {
	num := 64
	if len(n) > 0 {
		num = n[0]
	}
	b := make([]byte, num)
	for i, cache, remain := num-1, src.Int63(), numIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), numIdxMax
		}
		if idx := int(cache & numIdxMask); idx < len(numBytes) {
			b[i] = numBytes[idx]
			i--
		}
		cache >>= numIdxBits
		remain--
	}
	return Bytes2String(b)
}

// RandStrWithLower 随机生成小写字母字符串
func RandStrWithLower(n ...int) string {
	num := 64
	if len(n) > 0 {
		num = n[0]
	}
	b := make([]byte, num)
	for i, cache, remain := num-1, src.Int63(), lowerIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), lowerIdxMax
		}
		if idx := int(cache & lowerIdxMask); idx < len(lowerBytes) {
			b[i] = lowerBytes[idx]
			i--
		}
		cache >>= lowerIdxBits
		remain--
	}
	return Bytes2String(b)
}

// RandStrWithUpper 随机生成大写字母字符串
func RandStrWithUpper(n ...int) string {
	num := 64
	if len(n) > 0 {
		num = n[0]
	}
	b := make([]byte, num)
	for i, cache, remain := num-1, src.Int63(), upperIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), upperIdxMax
		}
		if idx := int(cache & upperIdxMask); idx < len(upperBytes) {
			b[i] = upperBytes[idx]
			i--
		}
		cache >>= upperIdxBits
		remain--
	}
	return Bytes2String(b)
}
