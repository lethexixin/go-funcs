package cryptos

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

import (
	"github.com/dgryski/go-skip32"
)

// EnSkip32 加密
func EnSkip32(key string, value uint32) (uint32, error) {
	s, err := skip32.New([]byte(key))
	if err != nil {
		return 0, err
	}
	return s.Obfus(value), nil
}

// DeSkip32 解密
func DeSkip32(key string, value uint32) (uint32, error) {
	s, err := skip32.New([]byte(key))
	if err != nil {
		return 0, err
	}
	return s.Unobfus(value), nil
}

// EncodeMD5 Md5加密(返回32位小写),该方法可以把字符串进行MD5加密,并把加密后的二进制转化为十六进制后的32位字符串返回
func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))
	return hex.EncodeToString(m.Sum(nil))
}

// EnSHA256 SHA256加密
func EnSHA256(value string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(value)))
}
