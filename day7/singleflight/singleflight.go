package singleflight

import "sync"

/**
 * user: ZY
 * Date: 2020/3/2 19:24
 */

//call 正在进行或已经结束的请求
type call struct{
	wg sync.WaitGroup
	val interface{}
	err error
}

//Group 主数据结构
type Group struct{
	mu sync.Mutex
	m map[string]*call
}

func (g *Group)Do(key string,fn func()(interface{},error))(interface{},error){
	g.mu.Lock()

	//延迟初始化
	if g.m == nil{
		g.m = make(map[string]*call)
	}

	//若获取到正在执行的请求则等待,且若请求结束,返回结果
	if c,ok:=g.m[key];ok{
		g.mu.Unlock()
		c.wg.Wait()
		return c.val,c.err
	}

	//若该请求没有放入Group中则生成call放入
	//group,call写加锁,且group写完释放锁
	c:=new(call)
	c.wg.Add(1)
	g.m[key]=c
	g.mu.Unlock()

	//调用一次fn(),并释放call所加的锁
	c.val,c.err=fn()
	c.wg.Done()

	//调用fn()结束该请求,并更新group
	//请求写锁
	g.mu.Lock()
	delete(g.m,key)
	g.mu.Unlock()

	//请求结束
	return c.val,c.err
}

