package dict

import (
	"errors"
	"fmt"
	"runtime"
	"testing"
)

func Test_dict(t *testing.T){
	defer func() {
		if e:=recover();e!=nil{
			fmt.Println(e)

			//打印调用栈信息
			buf := make([]byte, 2048)
			n := runtime.Stack(buf, false)
			stackInfo := fmt.Sprintf("%s", buf[:n])
			fmt.Println(stackInfo)
		}
	}()
	d:=GenerateDict(10)
	testUsed:=10
	for i:=0;i<testUsed;i++{
		d.Set(fmt.Sprintf("panyi-%d",i),fmt.Sprintf("value-%d",i))
	}

	for i:=0;i<testUsed;i++ {
		k,err:=d.Get(fmt.Sprintf("panyi-%d",i))
		if err!=nil{
			t.Error("fail",err)
		}
		if k!=fmt.Sprintf("value-%d",i){
			t.Error("fail",k)
		}
	}

	//fmt.Println(d.rehashSync(0))

	d.rehashAsyncInit(0)
	//var s0=d.hashLists[0].size
	//for i:=0;i<int(s0)+2;i++{
	//	d.rehashAsync()
	//}
	//PL(d.hashLists[0])

	for i:=0;i<testUsed;i++ {
		k,err:=d.Get(fmt.Sprintf("panyi-%d",i))
		if err!=nil{
			t.Error(fmt.Sprintf("fail %v,%s,%s",err,fmt.Sprintf("panyi-%d",i),k))
		}else if k!=fmt.Sprintf("value-%d",i){
			t.Error("fail",k)
		}
	}

	for i:=testUsed;i<testUsed+20;i++ {
		k,err:=d.Get(fmt.Sprintf("panyi-%d",i))
		if !errors.Is(err,Empty){
			t.Error(fmt.Sprintf("fail %v,%s,%s",err,fmt.Sprintf("panyi-%d",i),k))
		}else if k!=""{
			t.Error("fail",k)
		}
	}

	//PL(d.hashLists[0])
	//PL(d.hashLists[1])

	t.Log("pass")
}


func PL(dict *Hashlist){
	fmt.Println(fmt.Sprintf("----------------------used:%d,size:%d--------------------------",dict.used,dict.size))
	for i,v:=range dict.list{
		fmt.Printf("%d==>",i)
		vv:=v
		for {
			if vv==nil{
				break
			}else{
				fmt.Printf("		[%s]->[%s]",vv.key,vv.val)
				vv=vv.next
			}
		}
		fmt.Println()
	}
}

















