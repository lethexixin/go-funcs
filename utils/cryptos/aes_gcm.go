package cryptos

import (
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	cryptoRand "crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	mathRand "math/rand"
	"time"
)

const (
	gcmNonceSize = 12
)

var letters = []rune("1234567890abcdef")

func randSeq(n int) string {
	mathRand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[mathRand.Intn(len(letters))]
	}
	return string(b)
}

// aesGCMEncrypt aesGCM加密
func aesGCMEncrypt(plaintext, key, nonceSpec []byte) (ciphertext, nonce []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	if len(nonceSpec) != gcmNonceSize {
		nonce = make([]byte, gcmNonceSize)
		if _, err = io.ReadFull(cryptoRand.Reader, nonce); err != nil {
			return nil, nil, err
		}
	} else {
		nonce = nonceSpec
	}

	aesGcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	ciphertext = aesGcm.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nonce, nil
}

// aesGCMDecrypt aesGCM解密
func aesGCMDecrypt(cipherText, key, nonce []byte) (plaintext []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err = aesGcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// EncodeAESGCM aes-256-gcm AES/GCM/NoPadding
/*
AES_GCM加密解密的算法为: aes-256-gcm, nonce为随机12位字符串, padding为NoPadding, 常用于HTTP中, 过程如下:
1.客户端发送加密请求体 ---> 2.服务器接收请求体并解密 ---> 3.服务器处理请求,返回加密响应体 ---> 4.客户端接收响应体并解密

加密数据处理的流程为: 先对数据进行zlib 6 level压缩, 然后按照aes-256-gcm用key和nonce对数据进行加密, 最后对nonce拼接加密的数据进行Base64编码
即 AES_GCM加密: 压缩 --> AES加密 --> base64加密
*/

func EncodeAESGCM(data, key []byte) (base64ContentBytes []byte, nonceBytes []byte, err error) {
	if len(data) == 0 {
		return nil, nil, errors.New("length of data is zero")
	}

	// zlib compress
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	if _, err := w.Write(data); err != nil {
		return nil, nil, err
	}
	if err := w.Close(); err != nil {
		return nil, nil, err
	}
	data = b.Bytes()

	nonce := []byte(randSeq(gcmNonceSize))

	// encrypt
	cipherText, nonce, err := aesGCMEncrypt(data, key, nonce)
	if err != nil {
		return nil, nil, err
	}

	cipherText = append(nonce[:], cipherText[:]...)

	// base64 encode
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(cipherText)))
	base64.StdEncoding.Encode(buf, cipherText)
	return buf, nonce, nil
}

// DecodeAESGCM aes-256-gcm AES/GCM/NoPadding
/*
AES_GCM加密解密的算法为: aes-256-gcm, nonce为随机12位字符串, padding为NoPadding, 常用于HTTP中, 过程如下:
1.客户端发送加密请求体 ---> 2.服务器接收请求体并解密 ---> 3.服务器处理请求,返回加密响应体 ---> 4.客户端接收响应体并解密

解密数据处理的流程为: 先将nonce拼接加密的数据进行Base64解码, 然后按照aes-256-gcm用key和nonce对数据进行解密, 最后对数据进行zlib 6 level解压
即 AES_GCM解密: base64解密 --> AES解密 --> 解压缩
*/
func DecodeAESGCM(data, key []byte) (contentBytes []byte, err error) {
	if len(data) == 0 {
		return nil, errors.New("length of data is zero")
	}

	// base64 decode
	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(data))
	data, err = ioutil.ReadAll(decoder)
	if err != nil {
		return nil, err
	}
	if len(data) < gcmNonceSize {
		return nil, errors.New("length of data is wrong")
	}
	nonce := data[0:gcmNonceSize]
	// decrypt
	plaintext, err := aesGCMDecrypt(data[gcmNonceSize:], key, nonce)

	if err != nil {
		return nil, err
	}

	// zlib decompress
	r, err := zlib.NewReader(bytes.NewBuffer(plaintext))
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return content, nil
}
