package fastheap

import (
	//. "github.com/funny-falcon/go-fastheap"
	"testing"
	"sort"
	"fmt"
)
var _ = fmt.Print

type item struct {
	val int64
	ind int
}

func (i *item) Value() int64 {
	return i.val
}

func (i *item) Index() int {
	return i.ind
}

func (i *item) SetIndex(n int) {
	i.ind = n
}

func popall(h *IntHeap) []int {
	sl := make([]int, h.Size())
	for i, _ := range sl {
		v, _ := h.Pop()
		sl[i] = int(v.Value())
	}
	return sl
}

func testRangeIncrease(t *testing.T, from, to int) {
	for i:=from; i<to; i++ {
		h := IntHeap{}
		for j:=0; j<i; j++ {
			h.Insert(&item{val:int64(j)})
		}
		sl := popall(&h)
		if !sort.IntsAreSorted(sl) {
			t.Error(i, sl)
		}
	}
}

func testRangeDecrease(t *testing.T, from, to int) {
	for i:=from; i<to; i++ {
		h := IntHeap{}
		for j:=i; j>0; j-- {
			h.Insert(&item{val:int64(j)})
		}
		sl := popall(&h)
		if !sort.IntsAreSorted(sl) {
			t.Error(i, sl)
		}
	}
}

func TestInsertIncrease(t *testing.T) {
	testRangeIncrease(t, 1, 1024)
}

func TestInsertDecrease(t *testing.T) {
	testRangeDecrease(t, 1, 1024)
}
