package lists

import (
	"sort"
	"sync"
)

type Set struct {
	sync.RWMutex
	m map[interface{}]bool
}

// New 新建集合对象
func New(items ...interface{}) *Set {
	s := &Set{
		m: make(map[interface{}]bool, len(items)),
	}
	s.Add(items...)
	return s
}

// Add 添加元素
func (s *Set) Add(items ...interface{}) {
	s.Lock()
	defer s.Unlock()
	for _, v := range items {
		s.m[v] = true
	}
}

// Remove 删除元素
func (s *Set) Remove(items ...interface{}) {
	s.Lock()
	defer s.Unlock()
	for _, v := range items {
		delete(s.m, v)
	}
}

// Has 判断元素是否存在
func (s *Set) Has(items ...interface{}) bool {
	s.RLock()
	defer s.RUnlock()
	for _, v := range items {
		if _, ok := s.m[v]; !ok {
			return false
		}
	}
	return true
}

// Count 元素个数
func (s *Set) Count() int {
	return len(s.m)
}

// Clear 清空集合
func (s *Set) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = map[interface{}]bool{}
}

// Empty 空集合判断
func (s *Set) Empty() bool {
	return len(s.m) == 0
}

// List 无序列表
func (s *Set) List() []interface{} {
	s.RLock()
	defer s.RUnlock()
	list := make([]interface{}, 0, len(s.m))
	for item := range s.m {
		list = append(list, item)
	}
	return list
}

// SortList 排序列表,请参照此sort.Slice的写法,直接调用此方法无效
func (s *Set) SortList(pleaseDoNotUseThisFuncButYouCanSeeHowToUseInThisFunc chan error) []interface{} {
	list := s.List()
	sort.Slice(list, func(i, j int) bool {
		// TODO
		// 真正的写法应该是下面这个
		// return list[i].(People).Id > list[j].(People).Id
		return false
	})
	return list
}

// Union 并集
func (s *Set) Union(sets ...*Set) *Set {
	r := New(s.List()...)
	for _, set := range sets {
		for e := range set.m {
			r.m[e] = true
		}
	}
	return r
}

// DifferenceSet 差集
func (s *Set) DifferenceSet(sets ...*Set) *Set {
	r := New(s.List()...)
	for _, set := range sets {
		for e := range set.m {
			if _, ok := s.m[e]; ok {
				delete(r.m, e)
			}
		}
	}
	return r
}

// Intersection 交集
func (s *Set) Intersection(sets ...*Set) *Set {
	r := New(s.List()...)
	for _, set := range sets {
		for e := range s.m {
			if _, ok := set.m[e]; !ok {
				delete(r.m, e)
			}
		}
	}
	return r
}

// Complement 补集
func (s *Set) Complement(full *Set) *Set {
	r := New()
	for e := range full.m {
		if _, ok := s.m[e]; !ok {
			r.Add(e)
		}
	}
	return r
}
