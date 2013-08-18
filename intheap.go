package fastheap

type intItem struct {
	ref   IntValue
	value int64
}

// IntHeap is a heap holding int64 bit keys.
// It is minheap by default, and should be set as IntHeap{Max: true} if you want max heap
// I trust you will not change Max on non-empty heap :)
type IntHeap struct {
	Max  bool
	heap []*[pageSize]intItem
	size int
}

// Size returns size of a heap
func (h *IntHeap) Size() int {
	if h.size >= 3 {
		return h.size - 3
	}
	return 0
}

// Top returns top item and key value of a heap (minimum if Max == false, maximum otherwise
func (h *IntHeap) Top() (IntValue, int64) {
	if h.size > 3 {
		t := h.heap[0][3]
		return t.ref, t.value
	}
	return nil, 0
}

// Insert puts item into heap, preserving heap invariants
// Returns:
//   false, err - if error were encountered (inserting item has dirty index no)
//   false, nil - element were inserted and it's position is not at top
//   true, nil  - element were inserted at top position
func (h *IntHeap) Insert(tm IntValue) (bool, error) {
	if tm.Index() != 0 {
		return false, ErrInsert
	}

	h.ensureRoom()
	h.set(h.size, intItem{ref: tm, value: tm.Value()})
	h.size++
	ind := h.up(h.size - 1)
	return ind == 3, nil
}

// Remove removes item from heap, preserving heap invariants
// Returns:
//   false, err - if error were encountered (removing item has dirty index no)
//   false, nil - element were removed and it's position were not at top
//   true, nil  - element were removed from top position
func (h *IntHeap) Remove(tm IntValue) (bool, error) {
	if tm.Index() < 3 || tm.Index() >= h.size {
		return false, ErrRemove
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
	return i == 3, nil
}

// Reset fixes position of element, if it's key value were changed
func (h *IntHeap) Reset(tm IntValue) error {
	if tm.Index() < 3 || tm.Index() >= h.size {
		return ErrRemove
	}

	i := tm.Index()
	h.down(i)
	h.up(i)
	return nil
}

// Pop removes top item from heap (minimum if Max == false, maximum otherwise)
func (h *IntHeap) Pop() (IntValue, error) {
	if h.size <= 3 {
		return nil, ErrPop
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

// PopOrTop accepts cut mark and returns top item and key value of a heap, and were item popped cause it passed the cut mark
func (h *IntHeap) PopOrTop(cut int64) (IntValue, int64, bool) {
	if h.size > 3 {
		t := h.heap[0][3]
		if t.value > cut == h.Max {
			h.Pop()
			return t.ref, t.value, true
		}
		return t.ref, t.value, false
	}
	return nil, 0, false
}

func (h *IntHeap) up(j int) int {
	item := h.get(j)
	i := j/4 + 2
	if i == 2 || h.getValue(i) > item.value == h.Max {
		return j
	}
	h.move(i, j)
	j = i

	for {
		i = j/4 + 2
		if i == 2 || h.getValue(i) > item.value == h.Max {
			break
		}
		h.move(i, j)
		j = i
	}
	h.set(j, item)
	item.ref.SetIndex(j)
	return j
}

func (h *IntHeap) down(j int) {
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

func (h *IntHeap) downIndex(j int, e int64) int {
	last := h.size - 1
	if j > last/4+2 {
		return j
	}

	i1 := (j - 2) * 4
	i2 := i1 + 1
	i3 := i1 + 2
	i4 := i1 + 3
	chunk := h.getChunk(i1)

	e1 := chunk[i1&pageMask].value
	if i2 <= last {
		e2 := chunk[i2&pageMask].value
		if e2 > e1 == h.Max {
			i1 = i2
			e1 = e2
		}
	}
	if i3 <= last {
		e3 := chunk[i3&pageMask].value
		if i4 <= last {
			e4 := chunk[i4&pageMask].value
			if e4 > e3 == h.Max {
				i3 = i4
				e3 = e4
			}
		}
		if e3 > e1 == h.Max {
			i1 = i3
			e1 = e3
		}
	}

	if e1 > e == h.Max {
		j = i1
		e = e1
	}

	return j
}

func (h *IntHeap) ensureRoom() {
	if h.size > 0 {
		if h.size&pageMask == 0 && h.size>>pageLog == len(h.heap) {
			h.heap = append(h.heap, &[pageSize]intItem{})
		}
	} else {
		// initialization
		h.heap = make([]*[pageSize]intItem, 1)
		h.heap[0] = &[pageSize]intItem{}
		h.size = 3
	}
}

func (h *IntHeap) chomp() {
	chunks := ((h.size - 1) >> pageLog) + 1
	h.heap[h.size>>pageLog][h.size&pageMask] = intItem{}
	if chunks+1 < len(h.heap) {
		h.heap[len(h.heap)-1] = nil
		h.heap = h.heap[:chunks+1]
	}
}

func (h *IntHeap) get(i int) intItem {
	return h.heap[i>>pageLog][i&pageMask]
}

func (h *IntHeap) getChunk(i int) *[pageSize]intItem {
	return h.heap[i>>pageLog]
}

func (h *IntHeap) getValue(i int) int64 {
	return h.heap[i>>pageLog][i&pageMask].value
}

func (h *IntHeap) clear(i int) {
	h.heap[i>>pageLog][i&pageMask] = intItem{}
}

func (h *IntHeap) set(i int, item intItem) {
	h.heap[i>>pageLog][i&pageMask] = item
	item.ref.SetIndex(i)
}

func (h *IntHeap) move(from, to int) {
	item := h.heap[from>>pageLog][from&pageMask]
	h.heap[to>>pageLog][to&pageMask] = item
	item.ref.SetIndex(to)
}
