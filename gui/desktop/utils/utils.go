package utils

import (
	"crypto/aes"
	"crypto/cipher"
)

// AesDecryptByCBC
func AesDecryptByCBC(encrypted, key string) string {
	// 判断key长度
	keyLenMap := map[int]struct{}{16: {}, 24: {}, 32: {}}
	if _, ok := keyLenMap[len(key)]; !ok {
		panic("key长度必须是 16、24、32 其中一个")
	}
	// encrypted密文反解base64
	decodeString := []byte(encrypted)
	// key 转[]byte
	keyByte := []byte(key)
	// 创建一个cipher.Block接口。参数key为密钥，长度只能是16、24、32字节
	block, _ := aes.NewCipher([]byte(key))
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 选择加密模式
	blockMode := cipher.NewCBCDecrypter(block, keyByte[:blockSize])
	// 创建数组，存储解密结果
	decodeResult := make([]byte, blockSize)
	// 解密
	blockMode.CryptBlocks(decodeResult, decodeString)
	// 解码
	padding := PKCS7UNPadding(decodeResult)
	return string(padding)
}

// PKCS7UNPadding
func PKCS7UNPadding(originDataByte []byte) []byte {
	length := len(originDataByte)
	unpadding := int(originDataByte[length-1])
	return originDataByte[:(length - unpadding)]
}
