package fast

type fastCache struct {
	shards []*cacheShard
	shardMask uint64
	hash fnv64a
}

func NewFastCache(maxEntries,shardNum int, onEvicted func(key string,value interface{})) *fastCache{
	fastCache:=&fastCache{
		shards:    make([]*cacheShard,shardNum),
		shardMask: uint64(shardNum-1),
		hash:      newDefaultHasher(),
	}
	for i:=0;i<shardNum;i++{
		fastCache.shards[i] = newCacheShard(maxEntries,onEvicted)
	}
	return fastCache
}

func (c *fastCache) getShard(key string ) *cacheShard{
	hashKey:=c.hash.Sum64(key)
	return c.shards[hashKey&c.shardMask]
}

func (c *fastCache) Set(key string,value interface{}){
	c.getShard(key).set(key,value)
}

func (c *fastCache) Get(key string) interface{}{
	return c.getShard(key).get(key)
}

func (c *fastCache) Del(key string){
	c.getShard(key).del(key)
}

func ( c *fastCache) Len() int{
	length:=0
	for _,shard:=range c.shards{
		length+=shard.len()
	}
	return length
}

func (c *fastCache) DelOldest(){
	panic("no implements")
}