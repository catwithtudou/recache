package day6

import (
	"recache/day2/lru"
	"sync"
)

/**
 * user: ZY
 * Date: 2020/2/25 23:15
 */

type cache struct{
	mu sync.Mutex
	lru *lru.Cache
	cacheBytes int64
}

//add 向缓存中增加值
func (c *cache)add(key string,value ByteView){
	c.mu.Lock()
	defer c.mu.Unlock()

	//延迟初始化
	if c.lru==nil{
		c.lru = lru.New(c.cacheBytes,nil)
	}
	c.lru.Add(key,value)
}


//get 从缓存中取值
func (c *cache)get(key string)(value ByteView,ok bool){
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil{
		return
	}

	//若取到值则返回其value
	if v,ok:=c.lru.Get(key);ok{
		return v.(ByteView),ok
	}

	return
}