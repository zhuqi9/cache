package lru

import (
	"container/list"
	"github.com/zhuqi9/cache"
)

// lru是一个LRU cache。他不是并发安全的
type lru struct{
	// 缓存最大的容量
	// groupcache 使用的是最大存放entry 个数
	maxBytes int
	// 当一个entry 从缓存中移除时 调用该回调函数，默认为nil
	// groupcache 中的key 是任意的可比较类型； value 是 interface{}
	onEvicted func(key string,value interface{})

	//已经使用的字节数，只包括值，key不算
	usedBytes int

	ll *list.List
	cache map[string]*list.Element
}

type entry struct {
	key string
	value interface{}
}

func (e *entry)Len() int{
	return cache.CalcLen(e.value)
}

// New 创建一个 新的cache，如果 maxBytes 是0，表示 没有容量限制
func New (maxBytes int, onEvicted func(key string,value interface{})) cache.Cache{
	return &lru{
		maxBytes:maxBytes,
		onEvicted:onEvicted,
		ll:list.New(),
		cache:make(map[string]*list.Element),
	}
}

// Set 往cache 尾部增加一个元素，如果已经存在，则放入尾部，并更新值
func (l *lru) Set(key string, value interface{}){
	if e,ok:=l.cache[key];ok{
		l.ll.MoveToBack(e)
		en:=e.Value.(*entry)
		l.usedBytes = l.usedBytes -cache.CalcLen(en.value)+cache.CalcLen(value)
		en.value=value
		return
	}

	en:=&entry{key,value}
	e:=l.ll.PushBack(en )
	l.cache[key]=e

	l.usedBytes+=en.Len()
	if l.maxBytes>0&&l.usedBytes>l.maxBytes{
		l.DelOldest()
	}
}

// Get 从caceh中获取 key 对应的值，nil 表示key 不存在
func (l *lru) Get(key string) interface{}{
	if e,ok := l.cache[key];ok{
		l.ll.MoveToBack(e)
		return e.Value.(*entry).value
	}
	return nil
}

//Del 从cache 中删除key对应的元素
func (l *lru) Del(key string){
	if e,ok:=l.cache[key];ok{
		l.removeElement(e)
	}
}

//DelOldest 从cache 中删除最旧的记录
func (l *lru)DelOldest(){
	l.removeElement(l.ll.Front())
}

//Len 返回当前cache 中的记录数
func (l *lru)Len()int{
	return l.ll.Len()
}

func (l *lru) removeElement(e *list.Element){
	if e==nil{
		return
	}

	l.ll.Remove(e)
	en:=e.Value.(*entry)
	l.usedBytes-=en.Len()
	delete(l.cache,en.key)

	if l.onEvicted!=nil{
		l.onEvicted(en.key,en.value)
	}
}
