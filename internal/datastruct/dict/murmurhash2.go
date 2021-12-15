package dict

import "github.com/aviddiviner/go-murmur"

const SEED=000
func hash(b []byte)uint32{
	return murmur.MurmurHash2(b,SEED)
}