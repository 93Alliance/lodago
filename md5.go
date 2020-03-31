package lodago

import (
	"crypto/md5"
	"encoding/hex"
)

// Str2MD5 字符串转换成md5加密
func Str2MD5(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}
