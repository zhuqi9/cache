package fifo

import (
	"container/list"
	"github.com/zhuqi9/cache"
	)

//fifo 是一个FIFO cache，他不是并发安全的
type fifo struct {
	// 缓存最大的容量，单位字节
	// groupcache 使用的是最大存放 entry 个数
	maxBytes int
	// 当一个entry 从缓存中删除时调用该回调函数，默认为nil
	// groupcache 中的 key是任意的可比较类型： value是interface{}
	onEvicted func(key string ,value interface{})

	// 已使用的字节数，只包括值，key不断
	usedBytes int

	ll *list.List
	cache map[string]*list.Element
}

type entry struct {
	key string
	value interface{}
}

func (e *entry) Len() int {
	return cache.CalcLen(e.value)
}

//New 创建一个新的cache，如果 maxBytes 是0，表示没有容量限制
func New(maxBytes int, onEvicted func(key string,value interface{})) cache.Cache{
	return &fifo{
		maxBytes:maxBytes,
		onEvicted:onEvicted,
		ll: list.New(),
		cache:make(map[string]*list.Element),
	}
}

// Set 往 cache 尾部增加一个元素（如果存在则放入尾部，并修改）
func (f *fifo)Set(key string, value interface{}){
	if e,ok:=f.cache[key];ok{
		f.ll.MoveToBack(e)
		en:=e.Value.(*entry)
		f.usedBytes = f.usedBytes - cache.CalcLen(en.value)+ cache.CalcLen(value)
		en.value = value
		return
	}

	en:=&entry{key,value}
	e:=f.ll.PushBack(en )
	f.cache[key] = e

	f.usedBytes+=en.Len()
	if f.maxBytes>0&&f.usedBytes>f.maxBytes{
		f.DelOldest()
	}
}

// Get 从 cache中获取key 对应的值，nil表示key不存在
func (f *fifo)Get (key string) interface{}{
	if e,ok:=f.cache[key];ok{
		return e.Value.(*entry).value
	}
	return nil
}

// Del从cache 中删除key对应的记录
func (f *fifo)Del(key string){
	if e,ok:=f.cache[key];ok{
		f.removeElement(e)
	}
}

//DelOldest 从cache中删除最旧的记录
func (f *fifo)DelOldest(){
	f.removeElement(f.ll.Front())
}

func (f *fifo)Len() int {
	return f.ll.Len()
}

func (f *fifo)removeElement(e *list.Element){
	if e==nil{
		return
	}
	f.ll.Remove(e)
	en:=e.Value.(*entry)
	f.usedBytes-=en.Len()
	delete(f.cache,en.key)

	if f.onEvicted!=nil{
		f.onEvicted(en.key,en.value)
	}
}

