package dict

import (
	"fmt"
	"testing"
)

func Test_dict(t *testing.T){
	d:=GenerateDict(100000)
	testUsed:=1000
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

	fmt.Println(d.hashLists[0].used,d.hashLists[0].size)
	fmt.Println(d.rehashSync(0))
	fmt.Println(d.hashLists[0].used,d.hashLists[0].size)

	for i:=0;i<testUsed;i++ {
		k,err:=d.Get(fmt.Sprintf("panyi-%d",i))
		if err!=nil{
			t.Error("fail",err)
		}
		if k!=fmt.Sprintf("value-%d",i){
			t.Error("fail",k)
		}
	}


	t.Log("pass")
}


















