package main

/**
 * user: ZY
 * Date: 2020/2/27 9:53
 */

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	reCache "recache/day5"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *reCache.Group {
	return reCache.NewGroup("scores", 2<<10, reCache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

//startCacheServer 启动缓存服务器
func startCacheServer(addr string, addrs []string, re *reCache.Group) {
	peers := reCache.NewHTTPPool(addr)
	peers.Set(addrs...)

	re.RegisterPeers(peers)

	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))

}

//startAPIServer 启动一个 API 服务，与用户进行交互，用户感知
func startAPIServer(apiAddr string, re *reCache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := re.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())
		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}


func main(){
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "Geecache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	addrs := make([]string, 3)

	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	gee := createGroup()

	if api {
		go startAPIServer(apiAddr, gee)
	}

	startCacheServer(addrMap[port], []string(addrs), gee)
}