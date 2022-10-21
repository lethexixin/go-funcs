package cryptos

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"sync"
)

// GolangCharset 签名的字符编码类型
type GolangCharset string

// 字符编码类型常量
const (
	CharsetIso2022Jp         GolangCharset = "ISO-2022-JP"
	CharsetIso2022Cn                       = "ISO-2022-CN"
	CharsetIso2022Kr                       = "ISO-2022-KR"
	CharsetIso88595                        = "ISO-8859-5"
	CharsetIso88597                        = "ISO-8859-7"
	CharsetIso88598                        = "ISO-8859-8"
	CharsetBig5                            = "BIG5"
	CharsetGb18030                         = "GB18030"
	CharsetEucJp                           = "EUC-JP"
	CharsetEucKr                           = "EUC-KR"
	CharsetEucTw                           = "EUC-TW"
	CharsetShiftJis                        = "SHIFT_JIS"
	CharsetIbm855                          = "IBM855"
	CharsetIbm866                          = "IBM866"
	CharsetKoi8R                           = "KOI8-R"
	CharsetMacCyrillic                     = "x-mac-cyrillic"
	CharsetWindows1251                     = "WINDOWS-1251"
	CharsetWindows1252                     = "WINDOWS-1252"
	CharsetWindows1253                     = "WINDOWS-1253"
	CharsetWindows1255                     = "WINDOWS-1255"
	CharsetUtf8                            = "UTF-8"
	CharsetUtf16Be                         = "UTF-16BE"
	CharsetUtf16Le                         = "UTF-16LE"
	CharsetUtf32Be                         = "UTF-32BE"
	CharsetUtf32Le                         = "UTF-32LE"
	CharsetTis620                          = "WINDOWS-874"
	CharsetHzGb2312                        = "HZ-GB-2312"
	CharsetXIso10646Ucs43412               = "X-ISO-10646-UCS-4-3412"
	CharsetXIso10646Ucs42143               = "X-ISO-10646-UCS-4-2143"
)

// 当前类的指针
var sign *SignUtils

// 同步锁
var signOnce sync.Once

// SignUtils 签名类
type SignUtils struct {
	mapExtend *MapExtend
}

type MapExtend struct {
}

func (s *MapExtend) GetKeys(toSignMap *map[string]interface{}) (keys []string, err error) {
	for k := range *toSignMap {
		keys = append(keys, k)
	}
	return keys, err
}

// ToSignMap 形成符合签名格式的Map
func (s *MapExtend) ToSignMap(parameters *map[string]interface{}) map[string]interface{} {
	signMap := make(map[string]interface{})
	for k, v := range *parameters {
		if k != "extra" {
			signMap[k] = v
		} else {
			for eK, eV := range v.(map[string]interface{}) {
				signMap[eK] = eV
			}
		}
	}
	return signMap
}

// NewSign 实例化签名
func NewSign() *SignUtils {
	signOnce.Do(func() {
		sign = new(SignUtils)
		sign.mapExtend = new(MapExtend)
	})
	return sign
}

// GetUtf8Bytes 默认utf8字符串
func (s *SignUtils) GetUtf8Bytes(str string) []byte {
	b := []byte(str)
	return b
}

// SignTopRequest
/*
签名算法
parameters 要签名的数据项
secret 生成的publicKey
signMethod 签名的字符编码
*/
func (s *SignUtils) SignTopRequest(parameters *map[string]interface{}, secret string, signMethod GolangCharset) string {
	/**
	  1、第一步: 形成符合签名格式的Map
	  2、第二步: 按字典把Key的字母顺序排序
	  3、第三步: 把所有参数名和参数值串在一起
	  4、第四步: 使用MD5/HMAC加密,并把加密后的二进制转化为十六进制
	  5、第五步: 返回签名完成的大写的字符串
	*/

	//第一步: 形成符合签名格式的Map
	toSignMap := s.mapExtend.ToSignMap(parameters)
	//第二步: 按字典把Key的字母顺序排序
	keys, err := s.mapExtend.GetKeys(&toSignMap)
	if err != nil {
		return ""
	}
	sort.Strings(keys)
	//第三步: 把所有参数名和参数值串在一起
	var bb bytes.Buffer
	if CharsetUtf8 == signMethod {
		bb.WriteString(secret)
	}
	for i, v := range keys {
		_ = i
		val := (toSignMap)[v]
		if val != nil {
			bb.WriteString(v)
			bb.WriteString(fmt.Sprintf("%v", val))
		}
	}

	//第四步: 使用MD5/HMAC加密,并把加密后的二进制转化为十六进制
	var result string
	if CharsetUtf8 == signMethod {
		// cryptos.EncodeMD5() 方法可以把字符串进行MD5加密,并把加密后的二进制转化为十六进制后的字符串返回
		result = EncodeMD5(bb.String())
	} else {
		h := hmac.New(md5.New, s.GetUtf8Bytes(secret))
		h.Write(bb.Bytes())
		result = hex.EncodeToString(h.Sum(nil))
	}

	//第五步: 返回签名完成的大写的字符串
	return strings.ToUpper(result)
}
