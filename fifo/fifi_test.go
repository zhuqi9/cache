package fifo

import (
	"fmt"
	"github.com/matryer/is"
	//"github.com/zhuqi9/cache/fifo"
	"testing"
)


func TestSetGet(t *testing.T){
	is:=is.New(t)
	cache:=New(24,nil)
	cache.DelOldest()
	cache.Set("k1",1)
	v:=cache.Get("k1")
	is.Equal(v,1)

	cache.Del("k1")
	is.Equal(0,cache.Len())
}

func TestOnEvicted( t *testing.T){
	is:= is.New(t)
	keys:=make([]string,0,8)
	onEvicted:=func(key string,value interface{}){
		keys = append(keys,key)
	}
	cache:=New(16,onEvicted)
	res:=cache.(*fifo)
	cache.Set("k1",1)
	fmt.Println(res.cache["k1"].Value.(*entry).value)
	cache.Set("k2",2)
	fmt.Println(res)
	cache.Get("k1")
	cache.Set("k3",3)
	fmt.Println(res)
	cache.Get("k1")
	cache.Set("k4",4)
	fmt.Println(res)
	expected:=[]string{"k1","k2"}
	is.Equal(expected,keys)
	is.Equal(2,cache.Len())
}
