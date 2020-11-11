package lfu

import (
	"container/heap"
	"fmt"
	"github.com/zhuqi9/cache"
)

//lfu 是一个 LFU cache，他不是并发安全的

type lfu struct {
	// 缓存最大的容量，单位字节
	// groupcache 使用的是最大存放entry个数
	maxBytes int
	// 当一个entry 从缓存中移除时调用该回调函数，默认为nil
	//groupcache 中的key是任意课比较类型，value是interface{}
	onEvicted func(key string,value interface{})

	//已使用的字节数
	usedBytes int

	queue *queue
	cache map[string]*entry
}

//New 创建一个新的cache，如果maxBytes是0，表示没有容量限制
func New(maxBytes int, onEvicted func(key string,value interface{})) cache.Cache{
	q:=make(queue,0,1024)
	return &lfu{
		maxBytes:maxBytes,
		onEvicted:onEvicted,
		queue:&q,
		cache:make(map[string]*entry),
	}
}

//Set 往cache 增加一个元素（如果已经存在，则更新值，并增加权重，并重新构建堆）
func (l *lfu) Set (key string,value interface{}){
	if e,ok:=l.cache[key];ok{
		l.usedBytes = l.usedBytes - cache.CalcLen(e.value)+cache.CalcLen(value)
		l.queue.update(e,value,e.weight+1)
		return
	}

	en:=&entry{key:key,value:value}
	heap.Push(l.queue,en)
	l.cache[key]=en

	l.usedBytes+=en.Len()
	if l.maxBytes>0&&l.usedBytes>l.maxBytes{
		l.removeElement(heap.Pop(l.queue))
	}
}

//Get 从cache中获取key对应的value，nil表示key不存在
func (l *lfu)Get(key string)interface{}{
	if e,ok:=l.cache[key];ok{
		l.queue.update(e,e.value,e.weight+1)
		return e.value
	}
	return nil
}

// Del 从cache 中删除key对应的元素
func (l *lfu)Del(key string){
	if e,ok:=l.cache[key];ok{
		heap.Remove(l.queue,e.index)
		l.removeElement(e)
	}
}

//DelOldest 从cache总删除key对应的元素
func (l *lfu)DelOldest(){
	if l.queue.Len()==0{
		return
	}
	l.removeElement(heap.Pop(l.queue))
}

//Len 返回当前cache中的记录
func (l *lfu)Len() int{
	return l.queue.Len()
}

func (l *lfu)removeElement(x interface{}){
	if x==nil{
		return
	}
	en:=x.(*entry)
	fmt.Println("删除堆 ", en.key)
	delete(l.cache,en.key)
	l.usedBytes-=en.Len()
	if l.onEvicted!=nil{
		l.onEvicted(en.key,en.value)
	}
}
