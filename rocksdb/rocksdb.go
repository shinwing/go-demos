package main

import (
	"fmt"
	"strconv"
)

func main() {
	fmt.Println("start")

	db := NewRocksDB("/tmp/rocks/test/table")
	defer db.CloseDB()
	for i := 0; i < 100; i++ {
		keyStr := "hello" + strconv.Itoa(i)
		var key []byte = []byte(keyStr)
		db.Set(key, key)
		fmt.Println(i, keyStr)
	}

	it := db.GetIterator()
	defer it.Close()
	it.SeekToFirst()
	for ; it.Valid(); it.Next() {
		key := it.Key()
		fmt.Println("get data: ", string(key.Data()))
		defer key.Free()
	}

	fmt.Println("end")
}

//  CGO_CFLAGS="-I/root/workspace/golang/src/github.com/rocksdb/include" CGO_LDFLAGS="-L/root/workspace/golang/src/github.com/rocksdb -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy -llz4 " go build
