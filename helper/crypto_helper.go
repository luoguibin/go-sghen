package helper

import (
	"crypto/md5"
	"encoding/base64"
)

// MD5 ...
func MD5(word string) string {
	h := md5.New()
	h.Write([]byte(word))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// HmacMd5 ...验证
// func HmacMd5(src, key string) string {
// 	m := hmac.New(md5.New, []byte(key))
// 	m.Write([]byte(src))
// 	return base64.StdEncoding.EncodeToString(m.Sum(nil))
// }
