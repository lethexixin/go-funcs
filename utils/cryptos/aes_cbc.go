package cryptos

import (
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
)

// aesCBCEncrypt aesCBC加密
func aesCBCEncrypt(plaintext, key []byte, iv string) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// padding
	pad := byte(aes.BlockSize - (len(plaintext) % aes.BlockSize))
	plaintext = append(plaintext, bytes.Repeat([]byte{pad}, int(pad))...)

	var ciphertext []byte
	var ivByte []byte
	if len(iv) == 0 {
		ciphertext = make([]byte, aes.BlockSize+len(plaintext))
		ivByte = ciphertext[:aes.BlockSize]
		if _, err := io.ReadFull(rand.Reader, ivByte); err != nil {
			return nil, err
		}
		mode := cipher.NewCBCEncrypter(block, ivByte)
		mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)
	} else {
		ivByte = []byte(iv)
		ciphertext = make([]byte, len(plaintext))
		mode := cipher.NewCBCEncrypter(block, ivByte)
		mode.CryptBlocks(ciphertext, plaintext)
	}
	return ciphertext, nil
}

// aesCBCDecrypt aesCBC解密
func aesCBCDecrypt(cipherText, key []byte, iv string) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// The IV needs to be unique, but not secure. So it's common to
	// include it at the beginning of the ciphertext.
	if len(cipherText) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	var ivByte []byte
	if len(iv) == 0 {
		ivByte = cipherText[:aes.BlockSize]
		cipherText = cipherText[aes.BlockSize:]
	} else {
		ivByte = []byte(iv)
	}

	// CBC mode always works in whole blocks.
	if len(cipherText)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, ivByte)
	mode.CryptBlocks(cipherText, cipherText)
	bs := block.BlockSize()
	pad := int(cipherText[len(cipherText)-1])
	if pad > bs {
		return nil, errors.New("aes: invalid pkcs5 padding")
	}
	return cipherText[:len(cipherText)-pad], nil
}

// EncodeAESCBC aes-256-cbc AES/CBC/PKCS5Padding
/*
AES_CBC加密解密的算法为: aes-256-cbc, iv初始向量为随机16位字符串, padding为PKCS5Padding, 常用于HTTP中, 过程如下:
1.客户端发送加密请求体 ---> 2.服务器接收请求体并解密 ---> 3.服务器处理请求,返回加密响应体 ---> 4.客户端接收响应体并解密

加密数据处理的流程为: 先对数据进行zlib 6 level压缩, 然后按照aes-256-cbc用key和iv对数据进行加密, 最后对iv拼接加密的数据进行Base64编码
即 AES_CBC加密: 压缩 --> AES加密 --> base64加密
*/
func EncodeAESCBC(data, key []byte) (base64ContentBytes []byte, err error) {
	if len(data) == 0 {
		return nil, errors.New("length of data is zero")
	}

	// zlib compress
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}

	data = b.Bytes()
	// encrypt
	cipherText, err := aesCBCEncrypt(data, key, "")
	if err != nil {
		return nil, err
	}

	// base64 encode
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(cipherText)))
	base64.StdEncoding.Encode(buf, cipherText)
	return buf, nil
}

// DecodeAESCBC aes-256-cbc AES/CBC/PKCS5Padding
/*
AES_CBC加密解密的算法为: aes-256-cbc, iv初始向量为随机16位字符串, padding为PKCS5Padding, 常用于HTTP中, 过程如下:
1.客户端发送加密请求体 ---> 2.服务器接收请求体并解密 ---> 3.服务器处理请求,返回加密响应体 ---> 4.客户端接收响应体并解密

解密数据处理的流程为: 先将iv拼接加密的数据进行Base64解码, 然后按照aes-256-cbc用key和iv对数据进行解密, 最后对数据进行zlib 6 level解压
即 AES_CBC解密: base64解密 --> AES解密 --> 解压缩
*/
func DecodeAESCBC(data, key []byte) (contentBytes []byte, err error) {
	if len(data) == 0 {
		return nil, errors.New("length of data is zero")
	}

	// base64 decode
	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(data))
	data, err = ioutil.ReadAll(decoder)
	if err != nil {
		return nil, err
	}

	// aes decrypt
	plain, err := aesCBCDecrypt(data, key, "")
	if err != nil {
		return nil, err
	}

	// zlib decompress
	// r, err := zlib.NewReader(bytes.NewReader(plain))
	r, err := zlib.NewReader(bytes.NewBuffer(plain))
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return content, nil
}
