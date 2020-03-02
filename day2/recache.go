package day2

import (
	"fmt"
	"log"
	"sync"
)

/**
 * user: ZY
 * Date: 2020/2/25 23:22
 */

//Group 缓存命名空间
type Group struct{
	name string
	getter Getter
	mainCache cache
}

var (
	mu sync.RWMutex
	groups = make(map[string]*Group)
)


//NewGroup 新建一个缓存group
func NewGroup(name string,cacheBytes int64,getter Getter)*Group{
	if getter==nil{
		panic("nil Getter")
	}

	mu.Lock()
	defer mu.Unlock()

	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes:cacheBytes},
	}
	groups[name] = g

	return g
}

//GetGroup 从全局缓存groups中获取
func GetGroup(name string)*Group{
	//只读下用读锁
	mu.RLock()
	defer mu.RUnlock()
	g:=groups[name]
	return g
}


//Get 从缓存中取值
func (g *Group)Get(key string)(ByteView,error){
	if key==""{
		return ByteView{},fmt.Errorf("key os required")
	}

	//从实现的并发缓存中取值
	//若取值成功则返回
	if v,ok:=g.mainCache.get(key);ok{
		log.Println("[ReCache] hit")
		return v,nil
	}

	//若取值不成功则向源数据中加载
	return g.load(key)
}

//load 加载源数据
func (g *Group)load(key string)(value ByteView,err error){
	return g.getLocally(key)
}

//getLocally 调用回调函数获取源数据
func (g *Group)getLocally(key string)(ByteView,error){
	bytes,err:=g.getter.Get(key)
	if err!=nil{
		return ByteView{},err
	}

	value := ByteView{b:cloneBytes(bytes)}
	g.populateCache(key,value)
	return value,nil
}

//populateCache 将源数据添加到本地mainCache中
func (g *Group)populateCache(key string,value ByteView){
	g.mainCache.add(key,value)
}


//########################################################################


//Getter 解决缓存未命中时获取源数据的回调
type Getter interface {
	Get(key string) ([]byte,error) //回调函数
}


type GetterFunc func(key string)([]byte,error)

func (f GetterFunc)Get(key string)([]byte,error){
	return f(key)
}

