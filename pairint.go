package fastheap

type intPair struct {
	ref         IntValue
	value       int64
	left, right uint32
	up, down    uint32
}

type PairInt struct {
	Max  bool
	heap []*[pageSize]intPair
	size int
	len  uint32
	free uint32
}

func (h *PairInt) Size() int {
	return h.size
}

// Empty returns true when heap is empty
func (h *PairInt) Empty() bool {
	return h.size == 0
}

func (h *PairInt) Top() (IntValue, int64) {
	if h.size > 0 {
		root := &h.heap[0][0]
		t := h.at(root.down)
		return t.ref, t.value
	}
	return nil, 0
}

func (h *PairInt) Insert(tm IntValue) (bool, error) {
	if tm.Index() != 0 {
		return false, ErrInsert
	}

	t, i := h.getFree(tm)
	tm.SetIndex(int(i))
	h.size++
	root := &h.heap[0][0]
	top := h.at(root.down)
	if root.down == 0 {
		root.down = i
		t.left = i
		t.right = i
		return true, nil
	} else if t.value <= top.value {
		h.putUnder(i, root.down)
		root.down = i
		t.left = i
		t.right = i
		return true, nil
	} else {
		h.putUnder(root.down, i)
		return false, nil
	}
}

func (h *PairInt) Pop() (res IntValue, err error) {
	if h.size == 0 {
		return nil, ErrPop
	}
	root := &h.heap[0][0]
	top := h.at(root.down)
	res = top.ref
	res.SetIndex(0)
	ri := top.down
	*top = intPair{right: h.free}
	h.free = root.down
	h.size--
	if ri > 0 {
		r := h.at(ri)
		r.up = 0
		if r.left != ri {
			ri = r.left
			r = h.at(ri)
			for r.left != ri {
				li := r.left
				l := h.at(li)
				if l.value <= r.value {
					l.right = r.right
					h.at(r.right).left = li
					h.putUnder(li, ri)
					ri = li
					r = l
				} else {
					lli := l.left
					ll := h.at(lli)
					r.left = lli
					ll.right = ri
					h.putUnder(ri, li)
					ri = lli
					r = ll
				}
			}
		}
		root.down = ri
	} else {
		root.down = 0
	}
	return
}

func (h *PairInt) putUnder(up, down uint32) {
	tup := h.at(up)
	tdown := h.at(down)
	tdown.left = down
	tdown.right = down
	if tdown.up != 0 {
		tup.up = tdown.up
	}
	if tup.down == 0 {
		tdown.up = up
		tup.down = down
	} else {
		tfirst := h.at(tup.down)
		tlast := h.at(tfirst.left)
		tdown.right = tup.down
		tdown.left = tfirst.left
		tfirst.left = down
		tlast.right = down
		if tfirst.value >= tdown.value {
			tdown.up = up
			tup.down = down
			tfirst.up = 0
		} else {
			tdown.up = 0
		}
	}
}

func (h *PairInt) ati(i int) *intPair {
	return &h.heap[i>>pageLog][i&pageMask]
}

func (h *PairInt) at(i uint32) *intPair {
	return &h.heap[i>>pageLog][i&pageMask]
}

func (h *PairInt) getFree(tm IntValue) (t *intPair, i uint32) {
	if h.free != 0 {
		i = h.free
		t = h.at(h.free)
		h.free = t.right
	} else {
		if h.len > 0 {
			if h.len&pageMask == 0 && int(h.len>>pageLog) == len(h.heap) {
				h.heap = append(h.heap, &[pageSize]intPair{})
			}
		} else {
			// initialization
			h.heap = make([]*[pageSize]intPair, 1)
			h.heap[0] = &[pageSize]intPair{}
			h.len = 1
		}
		i = h.len
		h.len++
		t = h.at(i)
	}
	t.up = 0
	t.down = 0
	t.ref = tm
	t.value = tm.Value()
	return
}

func (h *PairInt) getChunk(i int) *[pageSize]intPair {
	return h.heap[i>>pageLog]
}

func (h *PairInt) getValue(i int) int64 {
	return h.heap[i>>pageLog][i&pageMask].value
}
