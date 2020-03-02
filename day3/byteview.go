package day3

/**
 * user: ZY
 * Date: 2020/2/25 23:10
 */

//ByteView 表示缓存值
type ByteView struct{
	b []byte
}


//Len 返回其所占内存大小
func (v ByteView)Len()int{
	return len(v.b)
}

//ByteSlice 返回一个拷贝,防止缓存值被外部程序修改
func (v ByteView)ByteSlice()[]byte{
	return cloneBytes(v.b)
}

//String 返回String类型的拷贝
func (v ByteView)String()string{
	return string(v.b)
}


//cloneBytes 拷贝缓存值
func cloneBytes(b []byte)[]byte{
	c:=make([]byte,len(b))
	copy(c,b)
	return c
}