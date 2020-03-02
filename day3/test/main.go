package main

import (
	"fmt"
	"log"
	"net/http"
	recache "recache/day3"
)

/**
 * user: ZY
 * Date: 2020/2/27 8:42
 */


var db = map[string]string{
	"Tom":  "119",
	"Jack": "121",
	"Sam":  "352",
}


func main(){
	recache.NewGroup("scores",2<<10,recache.GetterFunc(
		func(key string)([]byte,error) {
			log.Println("[SlowDB] search key",key)
			if v,ok:=db[key];ok{
				return []byte(v),nil
			}
			return nil, fmt.Errorf("%s not exist",key)
		}))

	addr:="localhost:8080"
	peers :=recache.NewHTTPPool(addr)
	log.Println("recache is running at",addr)
	log.Fatal(http.ListenAndServe(addr,peers))
}