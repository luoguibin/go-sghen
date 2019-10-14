package helper

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// MD5 ...
func MD5(word string) string {
	h := md5.New()
	h.Write([]byte(word))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// Sha256 ...
func Sha256(word string) string {
	h := sha256.New()
	h.Write([]byte(word))
	return hex.EncodeToString(h.Sum(nil))
	// ecab4881ee80ad3d76bb1da68387428ca752eb885e52621a3129dcf4d9bc4fd4 true
	// 295156ff7bd173eeb1f2c23a5219f6f38f350b5871f88fb364947ab083fa1f72
}

// HmacMd5 ...验证
// func HmacMd5(src, key string) string {
// 	m := hmac.New(md5.New, []byte(key))
// 	m.Write([]byte(src))
// 	return base64.StdEncoding.EncodeToString(m.Sum(nil))
// }
