package helper

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/base64"
)

// HmacMd5 ...验证
func HmacMd5(src, key string) string {
	m := hmac.New(md5.New, []byte(key))
	m.Write([]byte(src))
	return base64.StdEncoding.EncodeToString(m.Sum(nil))
}
