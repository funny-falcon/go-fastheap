package fastheap

import "fmt"

var _ = fmt.Errorf

type intArr struct {
	ref   IntValue
	value int64
	up    uint32
	down  [3]uint32
}

type ArrInt struct {
	Max  bool
	heap []*[pageSize]intArr
	size int
	top  uint32
	len  uint32
	free uint32
}

func (h *ArrInt) Size() int {
	return h.size
}

// Empty returns true when heap is empty
func (h *ArrInt) Empty() bool {
	return h.size == 0
}

func (h *ArrInt) Top() (IntValue, int64) {
	if h.size > 0 {
		t := h.at(h.top)
		return t.ref, t.value
	}
	return nil, 0
}

func (h *ArrInt) Insert(tm IntValue) (bool, error) {
	if tm.Index() != 0 {
		return false, ErrInsert
	}

	t, i := h.getFree(tm)
	tm.SetIndex(int(i))
	h.size++
	top := h.at(h.top)
	if h.top == 0 {
		h.top = i
		return true, nil
	} else if t.value <= top.value {
		h.putUnder(i, h.top)
		h.top = i
		return true, nil
	} else {
		h.putUnder(h.top, i)
		return false, nil
	}
}

func (h *ArrInt) Pop() (res IntValue, err error) {
	if h.size == 0 {
		return nil, ErrPop
	}
	top := h.at(h.top)
	res = top.ref
	res.SetIndex(0)
	curfree := h.free
	h.free = h.top
	h.size--
	if top.down[0] != 0 {
		if top.down[1] == 0 {
			h.top = top.down[0]
		} else if top.down[2] == 0 {
			if h.at(top.down[1]).value < h.at(top.down[0]).value {
				h.putUnder(top.down[1], top.down[0])
				h.top = top.down[1]
			} else {
				h.putUnder(top.down[0], top.down[1])
				h.top = top.down[0]
			}
		} else {
			h.sort2(top)
			h.putUnder(top.down[1], top.down[2])
			h.putUnder(top.down[0], top.down[1])
			h.top = top.down[0]
		}
	} else {
		h.top = 0
	}
	*top = intArr{up: curfree}
	return
}

func (h *ArrInt) putUnder(upi, downi uint32) {
	up := h.at(upi)
	down := h.at(downi)
	down.up = upi
	if up.down[2] != 0 {
		h.sort2(up)
		h.putUnder(up.down[1], up.down[2])
		h.putUnder(up.down[0], up.down[1])
		up.down[2] = 0
		up.down[1] = downi
	} else if up.down[1] != 0 {
		up.down[2] = downi
	} else if up.down[0] != 0 {
		up.down[1] = downi
	} else {
		up.down[0] = downi
	}
}

func (h *ArrInt) putUnderI(upi, downi uint32) {
	up := h.at(upi)
	h.at(downi).up = upi
	if up.down[2] != 0 {
		t := up
		ti := upi
		for {
			for t.down[2] != 0 {
				h.sort2(t)
				ti = t.down[1]
				t = h.at(ti)
			}
			p := h.at(t.up)
			if p.down[2] != 0 {
				pdi := p.down[2]
				p.down[2] = 0
				h.at(pdi).up = ti
				if t.down[0] == 0 {
					t.down[0] = pdi
				} else if t.down[1] == 0 {
					t.down[1] = pdi
				} else {
					t.down[2] = pdi
				}
				ti = p.down[0]
				t = h.at(ti)
			} else {
				pdi := p.down[1]
				p.down[1] = 0
				h.at(pdi).up = ti
				if t.down[0] == 0 {
					t.down[0] = pdi
				} else if t.down[1] == 0 {
					t.down[1] = pdi
				} else {
					t.down[2] = pdi
				}
				if p == up {
					break
				}
				ti = t.up
				t = p
			}
		}
		up.down[1] = downi
	} else if up.down[1] > 0 {
		up.down[2] = downi
	} else if up.down[0] > 0 {
		up.down[1] = downi
	} else {
		up.down[0] = downi
	}
}

func (h *ArrInt) sort2(it *intArr) {
	v0 := h.at(it.down[0]).value
	v1 := h.at(it.down[1]).value
	v2 := h.at(it.down[2]).value
	if v1 < v0 {
		it.down[0], it.down[1] = it.down[1], it.down[0]
		v1, v0 = v0, v1
	}
	if v2 < v1 {
		it.down[1], it.down[2] = it.down[2], it.down[1]
		v2, v1 = v1, v2
		if v1 < v0 {
			it.down[0], it.down[1] = it.down[1], it.down[0]
		}
	}
}

func (h *ArrInt) ati(i int) *intArr {
	return &h.heap[i>>pageLog][i&pageMask]
}

func (h *ArrInt) at(i uint32) *intArr {
	return &h.heap[i>>pageLog][i&pageMask]
}

func (h *ArrInt) valueAt(i uint32) int64 {
	return h.heap[i>>pageLog][i&pageMask].value
}

func (h *ArrInt) getFree(tm IntValue) (t *intArr, i uint32) {
	if h.free != 0 {
		i = h.free
		t = h.at(h.free)
		h.free = t.up
	} else {
		if h.len > 0 {
			if h.len&pageMask == 0 && int(h.len>>pageLog) == len(h.heap) {
				h.heap = append(h.heap, &[pageSize]intArr{})
			}
		} else {
			// initialization
			h.heap = make([]*[pageSize]intArr, 1)
			h.heap[0] = &[pageSize]intArr{}
			h.len = 1
		}
		i = h.len
		h.len++
		t = h.at(i)
	}
	t.up = 0
	t.down = [3]uint32{}
	t.ref = tm
	t.value = tm.Value()
	return
}

func (h *ArrInt) getChunk(i int) *[pageSize]intArr {
	return h.heap[i>>pageLog]
}
