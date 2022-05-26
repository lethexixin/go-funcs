package strs

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

import (
	uuid "github.com/satori/go.uuid"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// RandString 随机生成len位数的随机数,比如len=6则生成6位随机数,
// 需要配合init()方法中的rand.Seed(time.Now().Unix())使用
func RandString(len int) (s string) {
	for i := 0; i < len; i += 1 {
		a := rand.Int63n(10)
		s += fmt.Sprintf("%d", a)
	}
	return
}

// 生成随机UUID4
func UUID4() string {
	return uuid.NewV4().String()
}

// 判断传入的字符串str首或尾是否包含fixStr字符串,如果包含返回true
func StrPreAndSufIsContainFixStr(str string, fixStr string) bool {
	return strings.HasPrefix(str, fixStr) || strings.HasSuffix(str, fixStr)
}

// 压缩字符串,删除字符串中的多余空格,有多个空格时,仅保留一个空格,可用于将已格式化的json字符串压缩
func DelExtraSpace(src string) string {
	src = strings.ToLower(src) // 全部转小写
	// 删除字符串中的多余空格,有多个空格时,仅保留一个空格
	tempStr := strings.Replace(src, "  ", " ", -1) //替换tab为空格
	tempStr = strings.Replace(tempStr, "\r", " ", -1)
	tempStr = strings.Replace(tempStr, "\n", " ", -1)
	regStr := "\\s{2,}"                          //两个及两个以上空格的正则表达式
	reg, _ := regexp.Compile(regStr)             //编译正则表达式
	dst := make([]byte, len(tempStr))            //定义字符数组切片
	copy(dst, tempStr)                           //将字符串复制到切片
	spcIndex := reg.FindStringIndex(string(dst)) //在字符串中搜索
	for len(spcIndex) > 0 {                      //找到适配项
		dst = append(dst[:spcIndex[0]+1], dst[spcIndex[1]:]...) //删除多余空格
		spcIndex = reg.FindStringIndex(string(dst))             //继续在字符串中搜索
	}
	return strings.Trim(string(dst), " ")
}
