package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	tempStr := h.Sum(nil)
	return hex.EncodeToString(tempStr)
}
func MD5Encode(password, salt string) string {
	return Md5Encode(password + salt)
}
func Md5Decode(password, salt, code string) bool {
	return MD5Encode(password, salt) == code
}
