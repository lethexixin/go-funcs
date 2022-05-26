package slices

import (
	"reflect"
	"strings"
)

// 任意类型interface{} 转 []interface{}
func InterfaceToSlice(dataSlice interface{}) (result []interface{}) {
	if reflect.TypeOf(dataSlice).Kind() == reflect.Slice {
		data := reflect.ValueOf(dataSlice)
		for i := 0; i < data.Len(); i++ {
			ele := data.Index(i)
			result = append(result, ele.Interface())
		}
		return result
	}
	return nil
}

// 切片去重 空间换时间
func RemoveRepByMap(slc []interface{}) []interface{} {
	result := make([]interface{}, 0)
	tempMap := map[interface{}]byte{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}

// 切片去重 时间换空间
func RemoveRepByLoop(slc []interface{}) []interface{} {
	result := make([]interface{}, 0)
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i] == result[j] {
				flag = false // 存在重复元素，标识为false
				break
			}
		}
		if flag { // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	return result
}

// 通过反射的方式判断切片中是否有一样的内容,然后去重
func RemoveRepByReflect(src []interface{}) (dst []interface{}) {
	length := len(src)
	for i := 0; i < length; i++ {
		state := false
		for j := i + 1; j < length; j++ {
			if j > 0 && reflect.DeepEqual(src[i], src[j]) {
				state = true
				break
			}
		}
		if !state {
			dst = append(dst, src[i])
		}
	}
	return
}

func DropRepeat(list []string) []string {
	mapFlag := map[string]bool{}
	var rtn []string
	for _, i := range list {
		if i != "" && i != "0" {
			if _, ok := mapFlag[i]; !ok {
				rtn = append(rtn, i)
				mapFlag[i] = true
			}
		}
	}
	return rtn
}

func ListReduceList(list1, list2 []string) (list3 []string) {
	list1, list2 = DropRepeat(list1), DropRepeat(list2)
	s := "," + strings.Join(list2, ",") + ","
	for _, v := range list1 {
		if !strings.Contains(s, ","+v+",") {
			list3 = append(list3, v)
		}
	}
	return
}

func ListAddList(list1, list2 []string) []string {
	list1, list2 = DropRepeat(list1), DropRepeat(list2)
	for _, v := range list2 {
		s := "," + strings.Join(list1, ",") + ","
		if !strings.Contains(s, ","+v+",") {
			list1 = append(list1, v)
		}
	}
	return list1
}

func AddList(list, bid string) (s string) {
	if (list == "" || list == "[]") && bid != "" {
		s = "[" + bid + "]"
		return
	}
	if bid == "" {
		s = list
		return
	}
	flag := true
	list = list[1 : len(list)-1]
	list1 := strings.Split(list, ",")
	for _, v := range list1 {
		if v == bid {
			flag = false
		}
	}
	if flag {
		list1 = append(list1, bid)
	}
	s = "[" + strings.Join(list1, ",") + "]"
	return
}

func ReduceList(list, bid string) (s string) {
	if list == "" || list == "[]" {
		return
	}
	var list2 []string
	list = list[1 : len(list)-1]
	list1 := strings.Split(list, ",")
	for _, v := range list1 {
		if v != bid && v != "" && v != "0" {
			list2 = append(list2, v)
		}
	}
	if len(list2) == 0 {
		return
	}
	s = "[" + strings.Join(list2, ",") + "]"
	return
}

func ListReduce(list1, list2 []string) []string {
	if len(list1) == 0 || len(list2) == 0 {
		return list1
	}
	s := "," + strings.Join(list2, ",") + ","
	var list []string
	for _, v := range list1 {
		if !strings.Contains(s, ","+v+",") {
			list = append(list, v)
		}
	}
	return list
}

func IntListReduce(list1, list2 []int) []int {
	var newLs []int
	if len(list1) == 0 || len(list2) == 0 {
		return list1
	}
	for _, i := range list1 {
		end := false
		for _, j := range list2 {
			if i == j {
				end = true
				break
			}
		}
		if !end {
			newLs = append(newLs, i)
		}
	}
	return newLs
}

func Int64ListReduce(list1, list2 []int64) []int64 {
	var newLs []int64
	if len(list1) == 0 || len(list2) == 0 {
		return list1
	}
	for _, i := range list1 {
		end := false
		for _, j := range list2 {
			if i == j {
				end = true
				break
			}
		}
		if !end {
			newLs = append(newLs, i)
		}
	}
	return newLs
}

func Int64ListAppend(list []int64, i int64) []int64 {
	flag := false
	for _, v := range list {
		if v == i {
			flag = true
			break
		}
	}
	if !flag {
		list = append(list, i)
	}
	return list
}
