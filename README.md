# 并发容器
## List
### CopyOnWriteArrayList
#### 原理
    CopyOnWrite容器即写时复制的容器。
    通俗的理解是当我们往一个容器添加元素的时候，不直接往当前容器添加，而是先将当前容器进行Copy，复制出一个新的容器，然后新的容器里添加元素，添加完元素之后，再将原容器的引用指向新的容器。
    这样做的好处是我们可以对CopyOnWrite容器进行并发的读，而不需要加锁， 因为当前容器不会添加任何元素。所以CopyOnWrite容器也是一种读写分离的思想，读和写不同的容器。

#### 缺点
    内存占用问题
    因为CopyOnWrite的写时复制机制，所以在进行写操作的时候，内存里会同时驻扎两个对象的内存，旧的对象和新写入的对象（注意:在复制的时候只是复制容器里的引用，只是在写的时候会创建新对象添加到新容器里，而旧容器的对象还在使用，所以有两份对象内存）。如果这些对象占用的内存比较大，比如说200M左右，那么再写入100M数据进去，内存就会占用300M，那么这个时候很有可能造成频繁的Yong GC和Full GC。之前我们系统中使用了一个服务由于每晚使用CopyOnWrite机制更新大对象，造成了每晚15秒的Full GC，应用响应时间也随之变长。
    针对内存占用问题，可以通过压缩容器中的元素的方法来减少大对象的内存消耗，比如，如果元素全是10进制的数字，可以考虑把它压缩成36进制或64进制。或者不使用CopyOnWrite容器，而使用其他的并发容器，如ConcurrentHashMap。

    数据一致性问题
    CopyOnWrite容器只能保证数据的最终一致性，不能保证数据的实时一致性。所以如果你希望写入的的数据，马上能读到，请不要使用CopyOnWrite容器。
    关于C++的STL中，曾经也有过Copy-On-Write的玩法，参见陈皓的《C++ STL String类中的Copy-On-Write》，后来，因为有很多线程安全上的事，就被去掉了。

#### 方法
   + 并发安全方法
     + New(capacity int)
     + Add(ele element)
     + AddByPosition(index int, ele element)
     + Remove(index int)   
     + RemoveRange(fromIndex int, toIndex int)
     + Remove(index int)
   
   + 非并发安全方法
     + Get(index int)
     + Clear()
     + Size()
     + IsEmpty()
## Set

## Map
### SegmentMap
    分段map
### SkipListMap
    跳表map
### TreeMap
    红黑树map

