package list

import (
	"log"
	"runtime"
	"strconv"
	"sync"
	"testing"
)

var MAXPROCS int

func init() {
	MAXPROCS = runtime.NumCPU()
}

func TestCopyOnWriteList_Add(t *testing.T) {
	list := New(4)
	list.Add("1234")
	t.Logf("list: %v \n", list.elements)
}

func TestCopyOnWriteList_AddByPosition(t *testing.T) {
	cowl := New(4)
	//init data
	for i := 0; i < 10; i++ {
		cowl.Add("ele" + strconv.Itoa(i))
	}

	log.Printf("AddByPosition before Size:%d,elements:%v \n", cowl.Size(), cowl.elements)
	ret, err := cowl.AddByPosition(9, "ele-new")
	log.Printf("AddByPosition after Ret:%v ,Err:%v,Size:%d,elements:%v \n", ret, err, cowl.Size(), cowl.elements)
}

func TestCopyOnWriteList_Add2(t *testing.T) {
	cowl := New(4)
	wg := new(sync.WaitGroup)
	wg.Add(MAXPROCS + 1)
	for i := 0; i < MAXPROCS+1; i++ {
		go func() {
			for j := 0; j < 1000; j++ {
				ele := strconv.Itoa(i) + ":" + strconv.Itoa(j)
				cowl.Add(ele)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	log.Printf("MAXPROCS:%d,Size:%d,cowl:%v \n", MAXPROCS, cowl.Size(), cowl)
}

func TestCopyOnWriteList_Remove(t *testing.T) {
	cowl := New(4)
	//init data
	for i := 0; i < 10; i++ {
		cowl.Add("ele" + strconv.Itoa(i))
	}
	ele, err := cowl.Remove(9)
	log.Printf("removedElement:%v ,Err:%v,Size:%d,elements:%v \n", ele, err, cowl.Size(), cowl.elements)
}

func TestCopyOnWriteList_Remove2(t *testing.T) {
	cowl := New(4)
	//init data
	for i := 0; i < 10; i++ {
		cowl.Add("ele" + strconv.Itoa(i))
	}

	wg := new(sync.WaitGroup)
	wg.Add(4)
	for j := 0; j < 4; j++ {
		go func(index int) {
			log.Printf("Before removedIndex:%d, Size:%d,elements:%v \n", index, cowl.Size(), cowl.elements)
			ele, err := cowl.Remove(index)
			log.Printf("After removedIndex:%d, removedElement:%v ,Err:%v,Size:%d,elements:%v \n", index, ele, err, cowl.Size(), cowl.elements)
			//time.Sleep(time.Second * 2)
			wg.Done()
		}(j)
	}
	wg.Wait()
	log.Printf("Size:%d,elements:%v \n", cowl.Size(), cowl.elements)
}

func TestCopyOnWriteList_RemoveRange(t *testing.T) {
	cowl := New(4)
	//init data
	for i := 0; i < 20; i++ {
		cowl.Add("ele" + strconv.Itoa(i))
	}

	log.Printf("remove before Size:%d,elements:%v \n", cowl.Size(), cowl.elements)
	ret, err := cowl.RemoveRange(0, 1)
	log.Printf("remove after Ret:%v ,Err:%v,Size:%d,elements:%v \n", ret, err, cowl.Size(), cowl.elements)
}
