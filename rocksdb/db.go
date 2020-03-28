package main

import (
	"github.com/tecbot/gorocksdb"
)

func newDB(path string) *gorocksdb.DB {
	option := gorocksdb.NewDefaultOptions()
	option.SetCreateIfMissing(true)
	option.EnableStatistics()
	option.SetWriteBufferSize(8 * 1024)

	blockBasedTblOpt := gorocksdb.NewDefaultBlockBasedTableOptions()
	blockBasedTblOpt.SetBlockCache(gorocksdb.NewLRUCache(64 * 1024))
	blockBasedTblOpt.SetFilterPolicy(gorocksdb.NewBloomFilter(10))
	blockBasedTblOpt.SetIndexType(gorocksdb.KHashSearchIndexType)

	option.SetBlockBasedTableFactory(blockBasedTblOpt)
	option.SetAllowConcurrentMemtableWrites(false)

	store, err := gorocksdb.OpenDb(option, path)
	if err != nil {
		store.Close()
		panic(err)
	}
	return store
}

type RDB struct {
	*gorocksdb.DB
	writeOpts *gorocksdb.WriteOptions
	readOpts  *gorocksdb.ReadOptions
}

func NewRocksDb(path string) *RDB {
	db := &RDB{
		newDB(path),
		gorocksdb.NewDefaultReadOptions(),
		gorocksdb.NewDefaultWriteOptions(),
	}
	db.readOpts.SetFillCache(true)
	db.writeOpts.SetSync(true)
	return db
}

func (db *RDB) Set(key, value []byte) error {
	return db.Put(db.writeOpts, key, value)
}

func (db *RDB) Gett(key []byte) (*gorocksdb.Slice, error) {
	return db.Get(db.readOpts, key)
}

func (db *RDB) Close() {
	db.writeOpts.Destroy()
	db.readOpts.Destroy()
	db.Close()
}
