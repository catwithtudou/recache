package day6

/**
 * user: ZY
 * Date: 2020/2/27 9:21
 */



//PeerPicker
type PeerPicker interface {
	//根据传入的key选择相应节点
	PickPeer(key string)(peer PeerGetter,ok bool)
}

//PeerGetter
type PeerGetter interface {
	//用于从对应group中取值
	Get(group,key string)([]byte,error)
}

