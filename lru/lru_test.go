package lru

import (
	"fmt"
	"github.com/matryer/is"
	"testing"
)

func TestOnEvicted( t *testing.T){
	is := is.New(t)

	keys:=make([]string,0,8)
	onEvicted:=func(key string,value interface{}){
		keys = append(keys,key)
	}

	cache:=New(16,onEvicted)

	cache.Set("k1",1)
	fmt.Println(cache)
	cache.Set("k2",2)
	fmt.Println(cache)
	cache.Get("k1")
	cache.Set("k3",3)
	fmt.Println(cache)
	cache.Get("k1")
	cache.Set("k4",4)
	fmt.Println(cache)

	expected:=[]string{"k2","k3"}

	is.Equal(expected,keys)
	is.Equal(2,cache.Len())
}


