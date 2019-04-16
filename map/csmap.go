package _map

import "sync"

const (
	//默认初始化容量（segement的table）
	DEFAULT_INITIAL_CAPACITY int32 = 16

	//默认装载因子
	DEFAULT_LOAD_FACTOR float32 = 0.75

	//默认并发数（Segment的个数）
	DEFAULT_CONCURRENCY_LEVEL int32 = 16

	//段table的最大容量
	MAXIMUM_CAPACITY int32 = 1 << 30

	//段的最大数量
	MAX_SEGMENTS int32 = 1 << 16

	//每个segment中table的大小
	MIN_SEGMENT_TABLE_CAPACITY int32 = 2
)

type object = interface{}

type ConcurrentSegmentMap struct {
	segmentMask  int32
	segmentShift int32
	segments     []*Segment
}

func NewDefaultConcurrentSegmentMap() *ConcurrentSegmentMap {
	return NewConcurrentSegmentMap(DEFAULT_INITIAL_CAPACITY, DEFAULT_LOAD_FACTOR, DEFAULT_CONCURRENCY_LEVEL)
}

/**
	initialCapacity:整个Map的初始容量，实际操作的时候需要平均分给每个Segment
	loadFactor:负载因子,是给每个Segment内部使用的
	concurrentLevel:Segment的数量
 */
func NewConcurrentSegmentMap(initialCapacity int32, loadFactor float32, concurrentLevel int32) *ConcurrentSegmentMap {
	if initialCapacity > MAXIMUM_CAPACITY {
		initialCapacity = DEFAULT_INITIAL_CAPACITY
	}

	if loadFactor < 0 {
		loadFactor = DEFAULT_LOAD_FACTOR
	}

	if concurrentLevel > MAX_SEGMENTS {
		concurrentLevel = MAX_SEGMENTS
	}
	sshift := int32(0)
	ssise := int32(1)
	// 计算并行级别 ssize，因为要保持并行级别是 2 的 n 次方
	for ssise < concurrentLevel {
		ssise = ssise << 1
		sshift = sshift + 1
	}

	segmentShift := 32 - sshift
	segmentMask := ssise - 1

	//根据 initialCapacity 计算 Segment 数组中每个位置可以分到的大小
	c := initialCapacity / ssise
	//如果有余数则数量加1
	if (c * ssise) < initialCapacity {
		c = c + 1
	}

	tableCapacity := MIN_SEGMENT_TABLE_CAPACITY
	for tableCapacity < c {
		tableCapacity = tableCapacity << 1
	}

	//创建segment数组
	ss := make([]*Segment, ssise, ssise)
	//创建segment数组中的第一个元素s0
	threshold := float32(tableCapacity) * loadFactor
	table := make([]*Entry, tableCapacity, tableCapacity*2)
	s0 := newSegment(loadFactor, int32(threshold), table)
	ss = append(ss, s0)

	return &ConcurrentSegmentMap{
		segmentMask:  segmentMask,
		segmentShift: segmentShift,
		segments:     ss,
	}
}

func (csmap *ConcurrentSegmentMap) IsEmpty() bool {
	return true
}

func (csmap *ConcurrentSegmentMap) Size() int32 {
	return 0
}

func (csmap *ConcurrentSegmentMap) Get(key object) (object, error) {
	return nil, nil
}

func (csmap *ConcurrentSegmentMap) ContainsKey(key object) (bool, error) {
	return true, nil
}

func (csmap *ConcurrentSegmentMap) Put(key object, value object) (object, error) {
	return nil, nil
}

func (csmap *ConcurrentSegmentMap) PutIfAbsent(key object, value object) (object, error) {
	return nil, nil
}

//数据节点
type Entry struct {
	hash  uint32
	key   object
	value object
	next  *Entry
}

func (e *Entry) Key() object {
	return e.key
}

func (e *Entry) Value() object {
	return e.value
}

//段
type Segment struct {
	count      int32
	modCount   int32
	threshold  int32
	table      []*Entry
	loadFactor float32
	lock       *sync.Mutex
}

func newSegment(loadFactor float32, threshold int32, table []*Entry) *Segment {
	return &Segment{
		loadFactor: loadFactor,
		threshold:  threshold,
		table:      table,
	}
}

func (s *Segment) put(hash uint32, key object, value object) {

}
