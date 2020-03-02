package day5

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"recache/day5/consistenthash"
	"strings"
	"sync"
)

/**
 * user: ZY
 * Date: 2020/2/27 8:30
 */

const (
	defaultBasePath = "/_reCache/"
	defaultReplicas = 50
)

//HTTPPool 存储节点的HTTP连接池
type HTTPPool struct {
	self        string
	basePath    string
	mu          sync.Mutex
	peers       *consistenthash.Map
	httpGetters map[string]*httpGetter //映射远程节点与对应的httpGetter

}

//NewHTTPPool 初始化
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

//Set 实现一致性哈希算法,添加了传入的节点
func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.peers = consistenthash.New(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*httpGetter, len(peers))

	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseURL: peer + p.basePath}
	}
}

//PickPeer 根据Key选择节点返回节点对应的HTTP客户端
func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}
	return nil,false
}


var _ PeerPicker = (*HTTPPool)(nil)

//Log 服务端日志输出
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

//ServerHTTP 处理所有请求
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//判断是否有ReCache的默认路由前缀
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path : " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)

	// URL:/<basePath>/<groupName>/<key>
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusForbidden)
		return
	}

	groupName := parts[0]
	key := parts[1]

	//从缓存组里面取出缓存空间
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group:"+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	_, _ = w.Write(view.ByteSlice())

}

type httpGetter struct {
	baseURL string
}

func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	u := fmt.Sprintf("%v%v/%v", h.baseURL, url.QueryEscape(group), url.QueryEscape(key), )

	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned : %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body:%v", err)
	}
	return bytes, nil

}

var _ PeerGetter = (*httpGetter)(nil)
