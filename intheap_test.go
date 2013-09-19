package fastheap

import (
	//. "github.com/funny-falcon/go-fastheap"
	"fmt"
	"math/rand"
	"sort"
	"testing"
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

func popall(h IntInterface) []int {
	sl := make([]int, h.Size())
	for i, _ := range sl {
		v, _ := h.Pop()
		sl[i] = int(v.Value())
	}
	return sl
}

func _testing(t *testing.T, from, to int, fab func() IntInterface, fill func(h IntInterface, i int, t *testing.T)) {
	for i := from; i < to; i++ {
		h := fab()
		for k := 2; k > 0; k-- {
			fill(h, i, t)
			sl := popall(h)
			if len(sl) != i {
				t.Error("Popped less than inserted", h, sl)
			}
			if !sort.IntsAreSorted(sl) {
				t.Error("Not sorted", i, sl)
			}
			if !h.Empty() {
				t.Error("Not Empty", h)
			}
		}
	}
}

func popsome(h IntInterface, n int) {
	for i := 0; i < n; i++ {
		h.Pop()
	}
}

func _benching(N int, fab func() IntInterface, fill func(h IntInterface, i int)) {
	h := fab()
	k := 400
	for i := 0; i < k/4; i++ {
		fill(h, N/k)
		popsome(h, N/k/2)
	}
	for i := 0; i < k/4; i++ {
		fill(h, N/k)
		popsome(h, N/k/2)
	}
	for i := 0; i < k/4; i++ {
		fill(h, N/k)
		popsome(h, N/k/2*3)
	}
	for i := 0; i < k/4; i++ {
		fill(h, N/k)
		popsome(h, N/k/2*3)
	}
}

func intHeap() IntInterface {
	return &IntHeap{}
}

func pairInt() IntInterface {
	return &PairInt{}
}

func arrInt() IntInterface {
	return &ArrInt{}
}

func fillTestIncrease(h IntInterface, i int, t *testing.T) {
	for j := 0; j < i; j++ {
		/*
			if h.(*ArrInt).len != 0 {
				fmt.Printf("Before insert %d %+v %+v\n", j, h, h.(*ArrInt).heap[0])
			} else {
				fmt.Printf("Before insert %d %+v\n", j, h)
			}
		*/
		at_top, err := h.Insert(&item{val: int64(j)})
		//fmt.Printf("After insert %d %+v %+v\n", j, h, h.(*ArrInt).heap[0])
		if err != nil {
			t.Error(err)
		}
		if at_top && j != 0 {
			t.Error("Inserted at top", j)
		}
	}
}

func fillTestDecrease(h IntInterface, i int, t *testing.T) {
	for j := i; j > 0; j-- {
		at_top, err := h.Insert(&item{val: int64(j)})
		if err != nil {
			t.Error(err)
		}
		if !at_top {
			t.Error("Inserted not at top", j)
		}
	}
}

func fillTestRandom(h IntInterface, i int, t *testing.T) {
	for j := i; j > 0; j-- {
		_, err := h.Insert(&item{val: rand.Int63()})
		if err != nil {
			t.Error(err)
		}
	}
}

func fillBenchIncrease(h IntInterface, i int) {
	for j := 0; j < i; j++ {
		h.Insert(&item{val: int64(j)})
	}
}

func fillBenchDecrease(h IntInterface, i int) {
	for j := i; j > 0; j-- {
		h.Insert(&item{val: int64(j)})
	}
}

func fillBenchRandom(h IntInterface, i int) {
	for j := i; j > 0; j -= 2 {
		r := rand.Int63()
		r1 := rand.Int63()
		h.Insert(&item{val: r})
		h.Insert(&item{val: r1})
		/*
			h.Insert(&item{val: r1 + 1})
			h.Insert(&item{val: r1 + 2})
			h.Insert(&item{val: r1 + 3})
			h.Insert(&item{val: r + 2})
			h.Insert(&item{val: r + 3})
			h.Insert(&item{val: r + 4})
			h.Insert(&item{val: r + 5})
			h.Insert(&item{val: r + 7})
			h.Insert(&item{val: r + 6})
		*/
	}
}

func TestInsertIncrease(t *testing.T) {
	_testing(t, 1, 1024, intHeap, fillTestIncrease)
}

func TestInsertDecrease(t *testing.T) {
	_testing(t, 1, 1024, intHeap, fillTestDecrease)
}

func TestInsertRandom(t *testing.T) {
	_testing(t, 1, 1024, intHeap, fillTestRandom)
}

func BenchmarkInsertIncrease(b *testing.B) {
	_benching(b.N, intHeap, fillBenchIncrease)
}

func BenchmarkInsertDecrease(b *testing.B) {
	_benching(b.N, intHeap, fillBenchDecrease)
}

func BenchmarkInsertRandom(b *testing.B) {
	_benching(b.N, intHeap, fillBenchRandom)
}

const batch = 500000

func BenchmarkInsertIncrease1(b *testing.B) {
	for i := 0; i < b.N; i += batch {
		_benching(batch, intHeap, fillBenchIncrease)
	}
}

func BenchmarkInsertDecrease1(b *testing.B) {
	for i := 0; i < b.N; i += batch {
		_benching(batch, intHeap, fillBenchDecrease)
	}
}

func BenchmarkInsertRandom1(b *testing.B) {
	for i := 0; i < b.N; i += batch {
		_benching(batch, intHeap, fillBenchRandom)
	}
}

func TestInsertIncreasePair(t *testing.T) {
	_testing(t, 1, 1024, pairInt, fillTestIncrease)
}

func TestInsertDecreasePair(t *testing.T) {
	_testing(t, 1, 1024, pairInt, fillTestDecrease)
}

func TestInsertRandomPair(t *testing.T) {
	_testing(t, 1, 1024, pairInt, fillTestRandom)
}

func BenchmarkInsertIncreasePair(b *testing.B) {
	_benching(b.N, pairInt, fillBenchIncrease)
}

func BenchmarkInsertDecreasePair(b *testing.B) {
	_benching(b.N, pairInt, fillBenchDecrease)
}

func BenchmarkInsertRandomPair(b *testing.B) {
	_benching(b.N, pairInt, fillBenchRandom)
}

func BenchmarkInsertIncreasePair1(b *testing.B) {
	for i := 0; i < b.N; i += batch {
		_benching(batch, pairInt, fillBenchIncrease)
	}
}

func BenchmarkInsertDecreasePair1(b *testing.B) {
	for i := 0; i < b.N; i += batch {
		_benching(batch, pairInt, fillBenchDecrease)
	}
}

func BenchmarkInsertRandomPair1(b *testing.B) {
	for i := 0; i < b.N; i += batch {
		_benching(batch, pairInt, fillBenchRandom)
	}
}

func TestInsertIncreaseArr(t *testing.T) {
	_testing(t, 1, 1024, arrInt, fillTestIncrease)
}

func TestInsertDecreaseArr(t *testing.T) {
	_testing(t, 1, 1024, arrInt, fillTestDecrease)
}

func TestInsertRandomArr(t *testing.T) {
	_testing(t, 1, 1024, arrInt, fillTestRandom)
}

func BenchmarkInsertIncreaseArr(b *testing.B) {
	_benching(b.N, arrInt, fillBenchIncrease)
}

func BenchmarkInsertDecreaseArr(b *testing.B) {
	_benching(b.N, arrInt, fillBenchDecrease)
}

func BenchmarkInsertRandomArr(b *testing.B) {
	_benching(b.N, arrInt, fillBenchRandom)
}

func BenchmarkInsertIncreaseArr1(b *testing.B) {
	for i := 0; i < b.N; i += batch {
		_benching(batch, arrInt, fillBenchIncrease)
	}
}

func BenchmarkInsertDecreaseArr1(b *testing.B) {
	for i := 0; i < b.N; i += batch {
		_benching(batch, arrInt, fillBenchDecrease)
	}
}

func BenchmarkInsertRandomArr1(b *testing.B) {
	for i := 0; i < b.N; i += batch {
		_benching(batch, arrInt, fillBenchRandom)
	}
}
