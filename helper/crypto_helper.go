package helper

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
)

// HmacSha512 ...验证
func HmacSha512(src, key string) string {
	m := hmac.New(sha512.New, []byte(key))
	m.Write([]byte(src))
	if len(hex.EncodeToString(m.Sum(nil))) > 0 {
		return src
	}
	return ""
}
