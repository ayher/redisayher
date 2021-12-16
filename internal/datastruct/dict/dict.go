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
	defer func() {
		self.rehashAsync()
	}()
	ki:=hash([]byte(key))
	var hashListsItem *Hashlist
	if self.reHashFlag!=-1{
		hashListsItem=self.hashLists[1]
	}else{
		hashListsItem=self.hashLists[0]
	}
	index:=ki%hashListsItem.size
	hashListsItem.used++
	if hashListsItem.list[index]==nil{
		hashListsItem.list[index]=&listNode{
			key:key,
			val: val,
			next: nil,
		}
	}else{
		oneNod:=&listNode{
			key:key,
			val: val,
			next: hashListsItem.list[index],
		}
		hashListsItem.list[index]=oneNod
	}
	return nil
}

var Empty=fmt.Errorf("empty")

func (self *Dict)Get(key string)(string,error){
	defer func() {
		self.rehashAsync()
	}()
	ki:=hash([]byte(key))

	for i:=0;i< len(self.hashLists);i++{
		hashListsItem:=self.hashLists[i]

		index:=ki%hashListsItem.size
		indexNode:=hashListsItem.list[index]
		for{
			if indexNode==nil{
				break
			}
			if indexNode.key==key{
				return indexNode.val,nil
			}
			indexNode=indexNode.next
		}
		if self.reHashFlag==-1{
			break
		}
	}
	return "",Empty
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

func (self *Dict)rehashAsyncInit(rt uint8){
	if self.reHashFlag!=-1{
		return
	}
	var newSize uint32
	if rt==0{
		newSize=self.getNewLen(self.hashLists[0].used*2)
	}else if rt==1{
		newSize=self.getNewLen(self.hashLists[0].used/2)
	}

	newH1:=make([]*listNode,newSize)
	self.hashLists[1]=&Hashlist{
		list: newH1,
		size: newSize,
		used: 0,
	}

	self.reHashFlag=0
}

func (self *Dict)rehashAsync(){
	if self.reHashFlag==-1{
		return
	}else if self.reHashFlag>=int32(self.hashLists[0].size){
		self.hashLists[0]=self.hashLists[1]
		self.hashLists[1]=nil
		self.reHashFlag=-1
		return
	}

	node0:=self.hashLists[0].list[self.reHashFlag]

	for {
		if node0==nil{
			break
		}else{
			kk,vv:=node0.key,node0.val

			nk:=hash([]byte(kk))
			index:=nk%self.hashLists[1].size

			node1:=self.hashLists[1].list[index]
			if node1==nil{
				self.hashLists[1].list[index]=&listNode{
					key: kk,
					val: vv,
					next: nil,
				}
			}else{
				newListNode:=&listNode{
					key: kk,
					val: vv,
					next: node1,
				}
				self.hashLists[1].list[index]=newListNode
			}
			node0=node0.next
		}
	}

	self.hashLists[0].list[self.reHashFlag]=nil
	self.hashLists[1].used++
	self.reHashFlag++
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

func GenerateDict(size uint32) Dict {
	l:=make([]*listNode,size)
	return Dict{
		hashLists: [2]*Hashlist{
			{list: l,size: size,used: 0},
			nil,
		},
		reHashFlag: -1,
	}
}












