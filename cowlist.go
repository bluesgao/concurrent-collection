package cowlist

import (
	"errors"
	"strings"
	"sync"
)

type element = interface{}

type CopyOnWriteList struct {
	mutex    sync.Mutex //互斥锁
	elements []element  //存放数据的数组
}

func New() *CopyOnWriteList {
	return new(CopyOnWriteList).Init()
}

func (cowlist *CopyOnWriteList) Init() *CopyOnWriteList {
	cowlist.elements = make([]element, 8, 16)
	return cowlist
}

func (cowlist *CopyOnWriteList) ToString() string {
	eles := cowlist.elements //副本
	len := len(eles)
	if len == 0 {
		return "[]"
	}

	var str strings.Builder
	str.WriteString("[")
	for i := 0; i < len; i++ {
		//str.WriteString(fmt)
		//todo
	}
	str.WriteString("]")
	return str.String()
}

func (cowlist *CopyOnWriteList) Get(index int) element {
	eles := cowlist.elements //副本
	return eles[index]
}

func (cowlist *CopyOnWriteList) Add(ele element) bool {
	cowlist.mutex.Lock()
	//先复制，再追加，最后替换
	//newElements := l.elements
	//copy(l.elements, newElements)
	cowlist.elements = append(cowlist.elements, ele)
	defer cowlist.mutex.Unlock()
	return true
}

func (cowlist *CopyOnWriteList) AddByPosition(index int, ele element) (bool, error) {

	cowlist.mutex.Lock()
	len := len(cowlist.elements)
	//index不在有效范围内
	if index > len || index < 0 {
		return false, errors.New("index out of bounds")
	}
	//需要移动的元素个数
	numMoved := len - index
	//根据index将原始数组分成2段
	new1 := cowlist.elements[0:numMoved]
	new2 := cowlist.elements[numMoved:len]
	newElements := append(append(new1, ele), new2)
	cowlist.elements = newElements

	defer cowlist.mutex.Unlock()
	return true, nil
}

func (cowlist *CopyOnWriteList) Remove(index int) (element, error) {
	cowlist.mutex.Lock()
	len := len(cowlist.elements)
	//index不在有效范围内
	if index > len || index < 0 {
		return nil, errors.New("index out of bounds")
	}
	//被删除元素
	removedElement := cowlist.Get(index)
	numMoved := len - index - 1
	new1 := cowlist.elements[0:numMoved]
	new2 := cowlist.elements[numMoved+1 : len]
	newElements := append(new1, new2)
	cowlist.elements = newElements

	defer cowlist.mutex.Unlock()
	return removedElement, nil
}

func (cowlist *CopyOnWriteList) RemoveRange(fromIndex int, toIndex int) (bool, error) {
	cowlist.mutex.Lock()
	len := len(cowlist.elements)
	//fromIndex与toIndex不在有效范围内
	if fromIndex < 0 || toIndex > len || toIndex < fromIndex {
		return false, errors.New("index out of bounds")
	}

	newLen := len - (toIndex - fromIndex)
	numMoved := len - toIndex
	//删除到数组尾部
	if numMoved == 0 {
		cowlist.elements = cowlist.elements[0:newLen]
	} else {
		new1 := cowlist.elements[0:fromIndex]
		new2 := cowlist.elements[toIndex:newLen]
		newElements := append(new1, new2)
		cowlist.elements = newElements
	}

	defer cowlist.mutex.Unlock()
	return true, nil
}

func (cowlist *CopyOnWriteList) Clear() {
	cowlist.mutex.Lock()
	cowlist.elements = cowlist.elements[0:0]
	defer cowlist.mutex.Unlock()
}

func (cowlist *CopyOnWriteList) Size() int {
	return len(cowlist.elements)
}

func (cowlist *CopyOnWriteList) IsEmpty() bool {
	return len(cowlist.elements) == 0
}
