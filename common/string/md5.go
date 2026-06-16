package string

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(data string) string {
	if "" == data {
		return ""
	}
	byteData := []byte(data)
	h := md5.New()
	h.Write(byteData)
	return hex.EncodeToString(h.Sum(nil))
}
