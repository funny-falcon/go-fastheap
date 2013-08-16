package fastheap

type uintItem struct {
	ref   UintValue
	value uint64
}

/* UintHeap is a heap holding int64 bit keys. */
/* It is minheap by default, and should be set as IntHeap{Max: true} if you want max heap */
/* I trust you will not change Max on non-empty heap :) */
type UintHeap struct {
	Max  bool
	heap []*[256]uintItem
	size int
}

/* Size returns size of a heap */
func (h *UintHeap) Size() int {
	if h.size >= 3 {
		return h.size - 3
	}
	return 0
}

/* Top returns top item and key value of a heap (minimum if Max == false, maximum otherwise */
func (h *UintHeap) Top() (UintValue, uint64, bool) {
	if h.size > 3 {
		t := h.heap[0][3]
		return t.ref, t.value, true
	}
	return nil, 0, false
}

/* Insert puts item into heap, preserving heap invariants */
func (h *UintHeap) Insert(tm UintValue) error {
	if tm.Index() != 0 {
		return InsertError
	}

	h.ensureRoom()
	h.set(h.size, uintItem{ref: tm, value: tm.Value()})
	h.size++
	h.up(h.size - 1)
	return nil
}

/* Remove removes item from heap, preserving heap invariants */
func (h *UintHeap) Remove(tm UintValue) error {
	if tm.Index() < 3 || tm.Index() >= h.size {
		return RemoveError
	}

	i := tm.Index()
	tm.SetIndex(0)
	h.size--
	l := h.size
	if i != l {
		h.move(l, i)
		h.down(i)
		h.up(i)
	}
	h.chomp()
	return nil
}

/* Reset fixes position of element, if it's key value were changed */
func (h *UintHeap) Reset(tm UintValue) error {
	if tm.Index() < 3 || tm.Index() >= h.size {
		return RemoveError
	}

	i := tm.Index()
	h.down(i)
	h.up(i)
	return nil
}

/* Pop removes top item from heap (minimum if Max == false, maximum otherwise) */
func (h *UintHeap) Pop() (UintValue, error) {
	if h.size <= 3 {
		return nil, PopError
	}

	h.size--
	l := h.size
	t := h.heap[0][3]
	t.ref.SetIndex(0)
	if l > 3 {
		h.move(l, 3)
		h.down(3)
	} else {
		h.size--
	}
	h.chomp()
	return t.ref, nil
}

func (h *UintHeap) up(j int) {
	item := h.get(j)
	i := j/4 + 2
	if i == 2 || h.getValue(i) > item.value == h.Max {
		return
	}
	h.move(i, j)
	j = i

	for {
		i = j/4 + 2
		if i == j || h.getValue(i) > item.value == h.Max {
			break
		}
		h.move(i, j)
		j = i
	}
	h.set(j, item)
	item.ref.SetIndex(j)
}

func (h *UintHeap) down(j int) {
	var i int

	item := h.get(j)

	if i = h.downIndex(j, item.value); i == j {
		return
	}
	h.move(i, j)
	j = i

	for {
		i = h.downIndex(j, item.value)
		if i == j {
			break
		}
		h.move(i, j)
		j = i
	}
	h.set(j, item)
	item.ref.SetIndex(j)
}

func (h *UintHeap) downIndex(j int, e uint64) int {
	last := h.size - 1
	if j > last/4+2 {
		return j
	}
	var j2 int
	var e1, e2 uint64

	i1 := (j - 2) * 4
	i2 := i1 + 1

	e1 = h.getValue(i1)
	if e1 > e == h.Max {
		j = i1
		e = e1
	}

	if i2 <= last {
		if i2+1 <= last {
			e21, e22 := h.getValue(i2), h.getValue(i2+1)
			if e21 > e22 == h.Max {
				j2 = i2
				e2 = e21
			} else {
				j2 = i2 + 1
				e2 = e22
			}
		} else {
			j2 = i2
			e2 = h.getValue(i2)
		}
		if e2 > e == h.Max {
			j = j2
			e = e2
		}
	}
	return j
}
func (h *UintHeap) ensureRoom() {
	if h.size > 0 {
		if h.size&0xff == 0 && h.size>>8 == len(h.heap) {
			h.heap = append(h.heap, &[256]uintItem{})
		}
	} else {
		/* initialization */
		h.heap = make([]*[256]uintItem, 1)
		h.heap[0] = &[256]uintItem{}
		h.size = 3
	}
}

func (h *UintHeap) chomp() {
	chunks := ((h.size - 1) >> 8) + 1
	if chunks+1 < len(h.heap) {
		h.heap[len(h.heap)-1] = nil
		h.heap = h.heap[:chunks+1]
	} else {
		h.heap[h.size>>8][h.size&0xff] = uintItem{}
	}
}

func (h *UintHeap) get(i int) uintItem {
	return h.heap[i>>8][i&0xff]
}

func (h *UintHeap) getValue(i int) uint64 {
	return h.heap[i>>8][i&0xff].value
}

func (h *UintHeap) clear(i int) {
	h.heap[i>>8][i&0xff] = uintItem{}
}

func (h *UintHeap) set(i int, item uintItem) {
	h.heap[i>>8][i&0xff] = item
	item.ref.SetIndex(i)
}

func (h *UintHeap) move(from, to int) {
	item := h.heap[from>>8][from&0xff]
	h.heap[to>>8][to&0xff] = item
	item.ref.SetIndex(to)
}
