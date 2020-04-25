package main

import (
	"fmt"
	"math/rand"
)

const (
	maxLevel int = 16
)

type DataStorager interface {
	Get() []byte
}

type Iterator interface {
	HasNext() bool
	Next() Iterator
}

type Element struct {
	Score   float64
	Value   DataStorager
	forward []*Element
}

func newElement(score float64, value DataStorager, level int) *Element {
	return &Element{
		Score:   score,
		Value:   value,
		forward: make([]*Element, level),
	}
}

func (this *Element) HasNext() bool {
	return this.forward[0] != nil
}

func (this *Element) Next() Iterator {
	if this != nil {
		return this.forward[0]
	}
	return nil
}

func (this *Element) Val() DataStorager {
	return this.Value
}

func (this *Element) Set(v DataStorager) {
	this.Value = v
}

type SkipList struct {
	header *Element
	len    int
	level  int
}

func NewSkipList() *SkipList {
	return &SkipList{
		header: &Element{
			forward: make([]*Element, maxLevel),
		},
	}
}

func generateLevel() int {
	lv := 1
	for rand.Uint64()%2 == 1 && lv < maxLevel {
		lv++
	}
	return lv
}

func (this *SkipList) Front() Iterator {
	return this.header.forward[0]
}

func (this *SkipList) Len() int {
	return this.len
}

func (this *SkipList) _search(score float64) (*Element, []*Element) {
	update := make([]*Element, maxLevel)
	e := this.header
	for lv := this.level - 1; lv >= 0; lv-- {
		for e.forward[lv] != nil && e.forward[lv].Score < score {
			e = e.forward[lv]
		}
		update[lv] = e
	}
	e = e.forward[0]
	if e != nil && e.Score == score {
		return e, update
	}
	return nil, update
}

func (this *SkipList) Search(score float64) (DataStorager, bool) {
	e, _ := this._search(score)
	if e != nil {
		return e.Val(), true
	}
	return nil, false
}

func (this *SkipList) Insert(score float64, val DataStorager) DataStorager {
	e, update := this._search(score)
	if e != nil && e.Score == score {
		e.Set(val)
		return e.Val()
	}

	newLv := generateLevel()
	if newLv > this.level {
		newLv = this.level + 1
		update[this.level] = this.header
		this.level = newLv
	}
	e = newElement(score, val, newLv)
	for i := 0; i < newLv; i++ {
		e.forward[i] = update[i].forward[i]
		update[i].forward[i] = e
	}
	this.len++
	return e.Val()
}

func (this *SkipList) Delete(score float64) bool {
	e, update := this._search(score)
	if e != nil && e.Score == score {
		for i := 0; i < this.level; i++ { // todo
			if update[i].forward[i] != e {
				return false
			}
			update[i].forward[i] = e.forward[i]
		}
		this.len--
	}
	return true
}

func (this *SkipList) Print() {
	fmt.Println("\n-- SkipList-------------------------------------------")
	for i := this.level - 1; i >= 0; i-- {
		fmt.Println("level:", i)
		for e := this.Front(); e.HasNext(); e = e.Next() {
			fmt.Printf("%f:%s. ", e.(*Element).Score, string(e.(*Element).Val().Get()))
		}
		fmt.Println("\n -----------------------------------------------------")
	}
	fmt.Println("Current MaxLevel:", this.level)
}

// func (this *SkipList) Len() int
// func (this *SkipList) Reset()
// func (this *SkipList) Delete(n int64, s string) bool
// func (this *SkipList) DeleteByIndex(n int) bool
// func (this *SkipList) Search(n int64, s string) *Element
// func (this *SkipList) SearchByIndex(n int) *Element
