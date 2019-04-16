package list

import (
	"errors"
	"log"
	"sync"
)

type object = interface{}

const DEFAULT_CAPACITY = int(8)

const INT_MAX = int(^uint(0) >> 1) //int最大值

type CopyOnWriteList struct {
	mutex    sync.Mutex //互斥锁
	elements []object  //存放数据的数组
}

func New(capacity int) *CopyOnWriteList {
	log.Printf("INT_MAX:%d \n", INT_MAX)
	if capacity <= 0 || capacity > INT_MAX {
		capacity = DEFAULT_CAPACITY
	}
	return new(CopyOnWriteList).init(capacity)
}

func (cowlist *CopyOnWriteList) init(capacity int) *CopyOnWriteList {
	cowlist.elements = make([]object, 0, capacity)
	return cowlist
}

/*非线程安全*/
func (cowlist *CopyOnWriteList) Get(index int) object {
	eles := cowlist.elements //副本
	return eles[index]
}

func (cowlist *CopyOnWriteList) Add(ele object) bool {
	cowlist.mutex.Lock()
	//先复制，再追加，最后替换
	//newElements := l.elements
	//copy(l.elements, newElements)
	cowlist.elements = append(cowlist.elements, ele)
	defer cowlist.mutex.Unlock()
	return true
}

func (cowlist *CopyOnWriteList) AddByPosition(index int, ele object) (bool, error) {

	cowlist.mutex.Lock()
	oldLen := len(cowlist.elements)
	//index不在有效范围内
	if index >= oldLen || index < 0 {
		return false, errors.New("index out of bounds")
	}

	var newElements []object

	//头插
	if index == 0 {
		newElements = append(newElements, ele)
		for i := 0; i < oldLen; i++ {
			newElements = append(newElements, cowlist.elements[i])
		}
	} else if index == oldLen-1 { //尾插
		newElements = append(cowlist.elements, ele)
	} else { //中间插
		//将原始数组按照index分割成两部分
		part1 := cowlist.elements[0 : index-1]
		part2 := cowlist.elements[index:oldLen]
		//将分割后的两部分连接起来
		for i := 0; i < len(part1); i++ {
			newElements = append(newElements, part1[i])
		}

		//插入新值
		newElements = append(newElements, ele)

		for i := 0; i < len(part2); i++ {
			newElements = append(newElements, part2[i])
		}
	}
	cowlist.elements = newElements
	defer cowlist.mutex.Unlock()
	return true, nil
}

func (cowlist *CopyOnWriteList) Remove(index int) (object, error) {
	cowlist.mutex.Lock()
	oldLen := len(cowlist.elements)
	//index不在有效范围内
	if index >= oldLen || index < 0 {
		return nil, errors.New("index out of bounds")
	}
	//被删除元素
	removedElement := cowlist.Get(index)
	//log.Printf("removedElement:%v \n", removedElement)
	//第一个元素
	if index == 0 {
		cowlist.elements = cowlist.elements[index+1 : oldLen]
	} else if index == oldLen-1 { //最后一个元素
		cowlist.elements = cowlist.elements[0 : oldLen-1]
	} else { //中间元素
		var newElements []object
		//将原始数组按照index分割成两部分
		part1 := cowlist.elements[0 : index-1]
		part2 := cowlist.elements[index:oldLen]
		//将分割后的两部分连接起来
		for i := 0; i < len(part1); i++ {
			newElements = append(newElements, part1[i])
		}
		for i := 0; i < len(part2); i++ {
			newElements = append(newElements, part2[i])
		}
		cowlist.elements = newElements
	}
	defer cowlist.mutex.Unlock()
	return removedElement, nil
}

func (cowlist *CopyOnWriteList) RemoveRange(fromIndex int, toIndex int) (bool, error) {
	cowlist.mutex.Lock()
	oldLen := len(cowlist.elements)
	//fromIndex与toIndex不在有效范围内
	if fromIndex < 0 || toIndex >= oldLen || toIndex < fromIndex {
		return false, errors.New("index out of bounds")
	}

	//从第一个元素开始删除几个元素
	if fromIndex == 0 && toIndex >= 0 {
		cowlist.elements = cowlist.elements[toIndex+1 : oldLen]
	} else if fromIndex >= 0 && toIndex == oldLen-1 { //删除最后几个元素
		cowlist.elements = cowlist.elements[0:fromIndex]
	} else { //删除中间几个元素
		var newElements, part1, part2 []object
		//将原始数组按照rangeindex分割成两部分
		if fromIndex == toIndex {
			part1 = cowlist.elements[0 : fromIndex-1]
			part2 = cowlist.elements[fromIndex:oldLen]
		} else {
			part1 = cowlist.elements[0 : fromIndex-1]
			part2 = cowlist.elements[toIndex+1 : oldLen]
		}

		//将分割后的两部分连接起来
		for i := 0; i < len(part1); i++ {
			newElements = append(newElements, part1[i])
		}
		for i := 0; i < len(part2); i++ {
			newElements = append(newElements, part2[i])
		}
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

/*非线程安全*/
func (cowlist *CopyOnWriteList) Size() int {
	return len(cowlist.elements)
}

/*非线程安全*/
func (cowlist *CopyOnWriteList) IsEmpty() bool {
	return len(cowlist.elements) == 0
}
