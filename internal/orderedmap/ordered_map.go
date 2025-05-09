package orderedmap

import "slices"

type OrderedMap[T any] struct {
	keys   []string
	values map[string]T
}

func New[T any]() *OrderedMap[T] {
	return &OrderedMap[T]{
		keys:   []string{},
		values: make(map[string]T),
	}
}

func (om *OrderedMap[T]) Set(key string, value T) {
	if _, exists := om.values[key]; !exists {
		om.keys = append(om.keys, key)
	}
	om.values[key] = value
}

func (om *OrderedMap[T]) Get(key string) (T, bool) {
	v, ok := om.values[key]
	return v, ok
}

func (om *OrderedMap[T]) Delete(key string) {
	if _, exists := om.values[key]; !exists {
		return
	}
	delete(om.values, key)
	for i, k := range om.keys {
		if k == key {
			om.keys = slices.Delete(om.keys, i, i+1)
			break
		}
	}
}

// 判断 key 是否存在
func (om *OrderedMap[T]) Contains(key string) bool {
	_, exists := om.values[key]
	return exists
}

func (om *OrderedMap[T]) Clear() {
	om.keys = []string{}
	om.values = make(map[string]T)
}

// 按顺序返回所有元素（切片）
func (om *OrderedMap[T]) All() []T {
	result := make([]T, 0, len(om.keys))
	for _, key := range om.keys {
		result = append(result, om.values[key])
	}
	return result
}

// 遍历：传入一个函数
func (om *OrderedMap[T]) ForEach(f func(key string, value T)) {
	for _, key := range om.keys {
		f(key, om.values[key])
	}
}
