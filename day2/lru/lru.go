package lru

import "container/list"

/**
 * user: ZY
 * Date: 2020/2/25 22:37
 */


type Cache struct{
	maxBytes int64
	hasBytes int64  //当前已使用的内存
	ll *list.List
	cache map[string]*list.Element
	//可选:在条目被清除时执行回调
	OnEvicted func(key string,value Value)
}

type entry struct{
	key string
	value Value
}


type Value interface {
	Len() int //返回值所占用的内存大小
}

func New(maxBytes int64, onEvicted func(string,Value)) *Cache{
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

//Get 得到map中对应的链表节点,返回其Value并移动到队尾
//队首队尾相对的,在这里约定front为队尾
func (c *Cache) Get(key string) (value Value,ok bool){
	if element,ok :=c.cache[key];ok{
		c.ll.MoveToFront(element)
		entry:=element.Value.(*entry)
		return entry.value,true
	}
	return
}

//RemoveOldest 移除最近最小访问的节点即队首
func (c *Cache)RemoveOldest(){
	element:=c.ll.Back()
	if element!=nil{
		c.ll.Remove(element)
		entry:= element.Value.(*entry)
		//删除map中存放的Key
		delete(c.cache,entry.key)
		//减去所占用的内存
		c.hasBytes -= int64(len(entry.key)) + int64(entry.value.Len())
		//若回调函数不为nil则执行
		if c.OnEvicted!=nil{
			c.OnEvicted(entry.key,entry.value)
		}
	}
}

//Add 增加节点入缓存中,若键存在则更新值
func (c *Cache) Add(key string,value Value){
	if element,ok:=c.cache[key];ok{
		c.ll.MoveToFront(element)
		entry:=element.Value.(*entry)
		//更新占用内存差值
		c.hasBytes += int64(value.Len()) - int64(entry.value.Len())
		//更新value
		entry.value = value
		return
	}

	//若缓存不存在,则将节点移向队尾
	element:=c.ll.PushFront(&entry{
		key:   key,
		value: value,
	})

	c.cache[key]=element

	c.hasBytes += int64(len(key)) + int64(value.Len())

	//若新增内存大于内存最大值,则淘汰队首节点直至内存小于最大值
	for c.maxBytes !=0 && c.maxBytes<c.hasBytes{
		c.RemoveOldest()
	}

	return
}

//Len 获取添加的数据
func (c *Cache) Len() int{
	return c.ll.Len()
}