package day2

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

/**
 * user: ZY
 * Date: 2020/2/25 23:26
 */


func TestGetter(t *testing.T){
	var f Getter =GetterFunc(func(key string) ([]byte,error){
		return []byte(key),nil
	})
	expect := []byte("key")
	if v,_ :=f.Get("key");!reflect.DeepEqual(v,expect){
		t.Errorf("callback failed")
	}
}

var db = map[string]string{
	"Tom":  "111",
	"Jack": "589",
	"Sam":  "567",
}

func TestGroup_Get(t *testing.T) {
	loadCounts:=make(map[string]int,len(db))
	re:=NewGroup("scores",2<<10,GetterFunc(
		func(key string)([]byte,error) {
			log.Println("[SlowDB] search key",key)
			if v,ok:=db[key];ok{
				if _,ok:=loadCounts[key];!ok{
					loadCounts[key]=0
				}
				loadCounts[key] += 1
				return []byte(v),nil
			}
			return nil,fmt.Errorf("%s not exist",key)
		}))

	for k,v:=range db{
		if view,err:=re.Get(k);err!=nil||view.String()!=v{
			t.Fatalf("failed to get value of Tom")
		}
		if _,err:=re.Get(k);err!=nil||loadCounts[k]>1{
			t.Fatalf("cache %s miss",k)
		}
	}

	if view,err:=re.Get("unknown");err==nil{
		t.Fatalf("the value of unknow should be empty,but %s got",view)
	}
}