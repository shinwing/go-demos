package main

import assert "github.com/arl/assertgo"

type storager string

func (this storager) Get() []byte {
	return []byte(this)
}

func main() {
	list := NewSkipList()
	assert.True(list.Len() == 0)

	list.Insert(1, storager("12"))
	list.Insert(1, storager("32"))
	assert.True(list.Len() == 1)

	s, ok := list.Search(1)
	assert.True(ok)
	assert.True(string(s.Get()) == "32")

	list.Insert(2, storager("22"))
	list.Insert(13, storager("33"))
	list.Insert(5, storager("55"))
	list.Print()
	assert.True(list.Len() == 4)
}
