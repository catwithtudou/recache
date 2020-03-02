package day3

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

/**
 * user: ZY
 * Date: 2020/2/27 8:30
 */


const defaultBasePath = "/_reCache/"

//HTTPPool 存储节点的HTTP连接池
type HTTPPool struct {
	self string
	basePath string
}

//NewHTTPPool 初始化
func NewHTTPPool(self string)*HTTPPool{
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

//Log 服务端日志输出
func (p *HTTPPool)Log(format string,v ...interface{}){
	log.Printf("[Server %s] %s",p.self,fmt.Sprintf(format,v...))
}

//ServerHTTP 处理所有请求
func (p *HTTPPool)ServeHTTP(w http.ResponseWriter,r *http.Request){
	//判断是否有ReCache的默认路由前缀
	if !strings.HasPrefix(r.URL.Path,p.basePath){
		panic("HTTPPool serving unexpected path : "+r.URL.Path)
	}
	p.Log("%s %s",r.Method,r.URL.Path)

	// URL:/<basePath>/<groupName>/<key>
	parts := strings.SplitN(r.URL.Path[len(p.basePath):],"/",2)
	if len(parts)!=2{
		http.Error(w,"bad request",http.StatusForbidden)
		return
	}

	groupName := parts[0]
	key := parts[1]

	//从缓存组里面取出缓存空间
	group := GetGroup(groupName)
	if group==nil{
		http.Error(w,"no such group:"+groupName,http.StatusNotFound)
		return
	}

	view,err:=group.Get(key)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type","application/octet-stream")
	_, _ = w.Write(view.ByteSlice())

}




