package string

import (
	"encoding/hex"
	"strings"
)

func ByteToHex(data []byte) string {
	return strings.ToUpper(hex.EncodeToString(data))
}

func HexToByte(hexStr string) []byte {
	bytes, _ := hex.DecodeString(hexStr)
	return bytes
}
