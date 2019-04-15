package concurrentCollection

import (
	"errors"
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

func (l *CopyOnWriteList) Init() *CopyOnWriteList {
	l.elements = make([]element, 8, 16)
	return l
}

func (l *CopyOnWriteList) Get(index int) element {
	return l.elements[index]
}

func (l *CopyOnWriteList) Add(ele element) bool {
	l.mutex.Lock()
	//先复制，再追加，最后替换
	//newElements := l.elements
	//copy(l.elements, newElements)
	l.elements = append(l.elements, ele)
	defer l.mutex.Unlock()
	return true
}

func (l *CopyOnWriteList) AddByPosition(index int, ele element) (bool, error) {

	l.mutex.Lock()
	len := len(l.elements)
	//index不在有效范围内
	if index > len || index < 0 {
		return false, errors.New("index out of bounds")
	}
	//需要移动的元素个数
	numMoved := len - index
	//根据index将原始数组分成2段
	new1 := l.elements[0:numMoved]
	new2 := l.elements[numMoved:len]
	newElements := append(append(new1, ele), new2)
	l.elements = newElements

	defer l.mutex.Unlock()
	return true, nil
}

func (l *CopyOnWriteList) remove(index int) (element, error) {
	l.mutex.Lock()
	len := len(l.elements)
	//index不在有效范围内
	if index > len || index < 0 {
		return nil, errors.New("index out of bounds")
	}
	//被删除元素
	removedElement := l.Get(index)
	numMoved := len - index - 1
	new1 := l.elements[0:numMoved]
	new2 := l.elements[numMoved+1 : len]
	newElements := append(new1, new2)
	l.elements = newElements

	defer l.mutex.Unlock()
	return removedElement, nil
}

func (l *CopyOnWriteList) removeRange(fromIndex int, toIndex int)(bool,error){

}
