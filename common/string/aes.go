package string

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

const (
	key = "Qm8rT9vXcL2nP5sD7fH1jK0wZyA6bN3m"
	iv  = "pL9xV2cD8sH1mQ7a"
)

// pkcs7Padding 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	//判断缺少几位长度。最少1，最多 blockSize
	padding := blockSize - len(data)%blockSize
	//补足位数。把切片[]byte{byte(padding)}复制padding个
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// pkcs7UnPadding 填充的反向操作
func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("Encrypted data error")
	}
	//获取填充的个数
	unPadding := int(data[length-1])
	return data[:(length - unPadding)], nil
}

// AesEncrypt 加密
func AesEncrypt(data []byte) ([]byte, error) {
	var err error
	defer func() {
		if pr := recover(); pr != nil {
			err = errors.New("Encrypt data error")
		}
	}()
	//创建加密实例
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	//判断加密快的大小
	blockSize := block.BlockSize()
	//填充
	encryptBytes := pkcs7Padding(data, blockSize)
	//初始化加密数据接收切片
	crypted := make([]byte, len(encryptBytes))
	//使用cbc加密模式
	blockMode := cipher.NewCBCEncrypter(block, []byte(iv))
	//执行加密
	blockMode.CryptBlocks(crypted, encryptBytes)
	return crypted, nil
}

// AesDecrypt 解密
func AesDecrypt(data []byte) (string, error) {
	var err error
	defer func() {
		if pr := recover(); pr != nil {
			err = errors.New("Decrypt data error")
		}
	}()
	//创建实例
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	//使用cbc
	blockMode := cipher.NewCBCDecrypter(block, []byte(iv))
	//初始化解密数据接收切片
	crypted := make([]byte, len(data))
	//执行解密
	blockMode.CryptBlocks(crypted, data)
	//去除填充
	crypted, err = pkcs7UnPadding(crypted)
	if err != nil {
		return "", err
	}
	return string(crypted), nil
}
