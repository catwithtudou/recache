package day7

import (
	"fmt"
	"log"
	"recache/day6/singleflight"
	pb "recache/day7/recachepb"
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
	peers  PeerPicker
	//合并请求group
	loader *singleflight.Group
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
		loader: &singleflight.Group{},
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
//保证并发场景下针对相同的key,load的过程仅调用一次
func (g *Group)load(key string)(value ByteView,err error){
	//将加载源数据处理函数放入合并请求的loader中
	view,err:=g.loader.Do(key, func() (interface{},error) {
		if g.peers!=nil{
			if peer,ok:=g.peers.PickPeer(key);ok{
				if value,err = g.getFromPeer(peer,key);err==nil{
					return value,nil
				}
				log.Println("[ReCache] Failed to get from peer ",err)
			}
		}
		return g.getLocally(key)
	})

	if err!=nil{
		return view.(ByteView),nil
	}

	return
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

//RegisterPeers 缓存空间中注册选择节点接口,将HTTPPool注入到Group中
func (g *Group)RegisterPeers(peers PeerPicker){
	if g.peers !=nil{
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

//getFromPeer 访问http节点得到缓存值
func (g *Group)getFromPeer(peer PeerGetter,key string)(ByteView,error){
	req := &pb.Request{
		Group:                g.name,
		Key:                  key,
	}
	res:=&pb.Response{}
	err:=peer.Get(req,res)
	if err!=nil{
		return ByteView{},err
	}
	return ByteView{b:res.Value},nil
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
//########################################################################


