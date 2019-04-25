package _map

import (
	"bytes"
	"errors"
	"sync"
)

type concurrentMap struct {
	m     map[interface{}]interface{}
	mutex sync.Mutex
}

func NewConcurrentMap() *concurrentMap {
	return &concurrentMap{
		m: make(map[interface{}]interface{}),
	}
}

func (cmap *concurrentMap) Get(key interface{}) (err error, value interface{}) {
	if key == nil {
		return errors.New("key is nil"), nil
	}
	cmap.mutex.Lock()
	defer cmap.mutex.Unlock()
	v, ok := cmap.m[key]
	if ok {
		return nil, v
	}
	return errors.New("key no fund"), nil
}

func (cmap *concurrentMap) Put(key, value interface{}) (err error) {
	if key == nil || value == nil {
		return errors.New("key or value is nil")
	}

	cmap.mutex.Lock()
	defer cmap.mutex.Unlock()
	cmap.m[key] = value
	return nil
}

func (cmap *concurrentMap) Remove(key interface{}) (err error, value interface{}) {
	if key == nil {
		return errors.New("key is nil"), nil
	}

	cmap.mutex.Lock()
	defer cmap.mutex.Unlock()
	v, ok := cmap.m[key]
	if ok {
		delete(cmap.m, key)
		return nil, v
	}
	return errors.New("key no found"), nil
}

func (cmap *concurrentMap) ContainsKey(key interface{}) (error, bool) {
	if key == nil {
		return errors.New("key is nil"), false
	}
	cmap.mutex.Lock()
	defer cmap.mutex.Unlock()
	_, ok := cmap.m[key]
	if ok {
		return nil, true
	}
	return errors.New("key no found"), false
}

func (cmap *concurrentMap) Size() int {
	cmap.mutex.Lock()
	defer cmap.mutex.Unlock()
	return len(cmap.m)
}

func (cmap *concurrentMap) ToString() string {
	var buf bytes.Buffer
	buf.WriteString("[")
	cmap.mutex.Lock()
	cmap.mutex.Unlock()

	for k, v := range cmap.m {
		buf.WriteString(k.(string))
		buf.WriteString(":")
		buf.WriteString(v.(string))
		buf.WriteString(",")
	}
	buf.WriteString("]")
	return buf.String()
}
