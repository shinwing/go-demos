package main

import (
	"fmt"
	"strconv"
)

func main() {
	fmt.Println("start")

	db := NewRocksDb("/tmp/rocks/test/table")
	defer db.Close()

	for i := 0; i < 10000; i++ {
		keyStr := "hello" + strconv.Itoa(i)
		var key []byte = []byte(keyStr)
		db.Set(key, key)
		fmt.Println(i, keyStr)
		slice, err := db.Gett(key)
		if err != nil {
			fmt.Println("error when get data：", key, err)
			continue
		}
		fmt.Println("get data：", slice.Size(), string(slice.Data()))
	}

	fmt.Println("end")
}
