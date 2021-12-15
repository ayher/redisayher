package dict

import "testing"

func Test_hash(t *testing.T){
	println(hash([]byte("foo")))       // == 1412061192
}
