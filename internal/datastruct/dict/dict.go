package dict

import "fmt"

type Dict struct {
	hashLists [2]*Hashlist
	reHashFlag int32
}

type Hashlist struct {
	list []*listNode
	used uint32
	size uint32
}

type listNode struct {
	key string
	val string
	next *listNode
}

func (self *Dict)Set(key,val string)error{
	ki:=hash([]byte(key))
	index:=ki%self.hashLists[0].size
	self.hashLists[0].used++
	if self.hashLists[0].list[index]==nil{
		self.hashLists[0].list[index]=&listNode{
			key:key,
			val: val,
			next: nil,
		}
	}else{
		oneNod:=&listNode{
			key:key,
			val: val,
			next: self.hashLists[0].list[index],
		}
		self.hashLists[0].list[index]=oneNod
	}
	return nil
}

func (self *Dict)Get(key string)(string,error){
	ki:=hash([]byte(key))
	index:=ki%self.hashLists[0].size
	indexNode:=self.hashLists[0].list[index]
	for{
		if indexNode==nil{
			return "",fmt.Errorf("empty")
		}
		if indexNode.key==key{
			return indexNode.val,nil
		}
		indexNode=indexNode.next
	}
	
}

func (self *Dict)rehashSync(rt uint8) int {
	if self.reHashFlag==-1{
		return 0
	}
	var newSize uint32
	if rt==0{
		newSize=self.getNewLen(self.hashLists[0].used*2)
	}else if rt==1{
		newSize=self.getNewLen(self.hashLists[0].used/2)
	}

	newH1:=make([]*listNode,newSize)

	newH1Used:=uint32(0)

	sameC:=0
	for _,item:=range self.hashLists[0].list{
		it:=item
		for{
			if it==nil{
				break
			}
			newH1Used++
			k:=it.key
			v:=it.val
			nk:=hash([]byte(k))
			index:=nk%newSize
			indexNode:=newH1[index]
			if indexNode==nil{
				newH1[index]=&listNode{
					key:k,
					val: v,
					next: nil,
				}
			}else{
				//same(indexNode.key,k,newSize)
				newindexNode:=&listNode{
					key:k,
					val: v,
					next: indexNode,
				}
				newH1[index]=newindexNode
				sameC++
			}
			it=it.next
		}
	}

	self.hashLists[0]=&Hashlist{
		list: newH1,
		size: newSize,
		used: newH1Used,
	}

	self.hashLists[1]=nil
	return sameC
}

//func same(k1,k2 string,size uint32){
//	nk1:=hash([]byte(k1))
//	index1:=nk1%size
//
//	nk2:=hash([]byte(k2))
//	index2:=nk1%size
//
//	fmt.Println(k1,nk1,index1,size)
//	fmt.Println(k2,nk2,index2,size)
//
//}

func (self *Dict)getNewLen(nlen uint32) uint32 {
	newL:=nlen*2
	var t=uint32(1)
	for{
		if t>=newL{
			return t
		}
		t=t*2
	}
}

func GenerateDict(size uint32) Dict{
	l:=make([]*listNode,size)
	return Dict{
		hashLists: [2]*Hashlist{
			{list: l,size: size,used: 0},
			nil,
		},
	}
}