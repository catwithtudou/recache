package day4

import (
	"hash/crc32"
	"sort"
	"strconv"
)

/**
 * user: ZY
 * Date: 2020/2/27 8:57
 */

//Hash 采取依赖注入的方式,允许用于替换成自定义的Hash函数,也方便测试时特换
type Hash func(data []byte) uint32

//Map 核心结构
type Map struct {
	hash     Hash
	replicas int            //虚拟节点倍数
	keys     []int          //哈希环
	hashMap  map[int]string //虚拟节点与真实节点的映射表
}

//New 初始化
func New(replicas int, fn Hash) *Map {
	m := &Map{
		hash:     fn,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}

	return m
}

//Add 添加真实节点/机器
func (m *Map) Add(keys ...string) {
	for _, v := range keys {
		//一个真实节点创造m.replicas的虚拟节点
		for i := 0; i < m.replicas; i++ {
			hash:=int(m.hash([]byte(strconv.Itoa(i)+v)))
			//fmt.Println(v+":"+strconv.Itoa(hash))
			m.keys = append(m.keys,hash)
			m.hashMap[hash]=v
		}
	}
	//环上哈希值排序
	sort.Ints(m.keys)
}

//Get 选择节点
func (m *Map)Get(key string)string{
	if len(m.keys)==0{
		return ""
	}

	hash := int(m.hash([]byte(key)))
	//顺时针找到第一个匹配的虚拟节点下标idx
	idx:=sort.Search(len(m.keys),func (i int)bool{
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]
}


