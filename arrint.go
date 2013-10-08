package fastheap

import "fmt"

var _ = fmt.Errorf

type intArr struct {
	ref        IntValue
	value      int64
	up         uint32
	upi, downl uint16
	down       [4]uint32
}

type ArrInt struct {
	Max  bool
	heap []*[pageSize]intArr
	size int
	top  uint32
	len  uint32
	free uint32
}

func (h *ArrInt) String() string {
	res := fmt.Sprintf("{s: %d t: %d l: %d f: %d [", h.size, h.top, h.len, h.free)
	for i, a := range h.heap {
		for j, ar := range *a {
			if ar.ref != nil {
				res += fmt.Sprintf("%d={v: %d u: %d ui: %d dl: %d [%d %d %d %d]} ", i*pageSize+j, ar.value, ar.up, ar.upi, ar.downl, ar.down[0], ar.down[1], ar.down[2], ar.down[3])
				for k, kk := range ar.down {
					if kk > 0 && h.at(kk).value < ar.value {
						panic(fmt.Errorf("child value less %d %d %d %+v %+v", i*pageSize+j, kk, k, ar, h.at(kk)))
					}
				}
			}
		}
	}
	res += "]}"
	return res
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
	res = h.at(h.top).ref
	h.delete(h.top)
	return res, nil
}
func (h *ArrInt) Remove(tm IntValue) (bool, error) {
	i := uint32(tm.Index())
	if i < 1 || i > h.len {
		return false, ErrRemove
	}
	attop := i == h.top
	h.delete(i)
	return attop, nil
}

func (h *ArrInt) delete(at uint32) {
	t := h.at(at)
	t.ref.SetIndex(0)
	curfree := h.free
	h.free = at
	h.size--
	var downi uint32
	switch t.downl {
	case 0:
		downi = 0
	case 1:
		downi = t.down[0]
	case 2:
		if h.at(t.down[1]).value < h.at(t.down[0]).value {
			h.putUnder(t.down[1], t.down[0])
			downi = t.down[1]
		} else {
			h.putUnder(t.down[0], t.down[1])
			downi = t.down[0]
		}
	case 3:
		h.sort3(&t.down)
		downi = t.down[0]
	case 4:
		h.sort4(&t.down)
		downi = t.down[0]
	}
	if at == h.top {
		h.top = downi
		down := h.at(downi)
		down.up = 0
		down.upi = 0
	} else if downi == 0 {
		up := h.at(t.up)
		switch t.upi {
		case 0:
			up.down[0] = up.down[1]
			h.at(up.down[1]).upi--
			fallthrough
		case 1:
			up.down[1] = up.down[2]
			h.at(up.down[2]).upi--
			fallthrough
		case 2:
			up.down[2] = up.down[3]
			h.at(up.down[3]).upi--
			fallthrough
		case 3:
			up.down[3] = 0
		}
		up.downl--
	} else {
		up := h.at(t.up)
		down := h.at(downi)
		down.up = t.up
		down.upi = t.upi
		up.down[t.upi] = downi
	}
	*t = intArr{up: curfree}
	return
}

func (h *ArrInt) putUnder(upi, downi uint32) {
	up := h.at(upi)
	down := h.at(downi)
	down.up = upi
	if up.downl == 4 {
		h.sort4(&up.down)
		up.down[3] = 0
		up.down[2] = 0
		up.down[1] = downi
		h.at(up.down[0]).upi = 0
		down.upi = 1
		up.downl = 2
	} else {
		up.down[up.downl] = downi
		down.upi = up.downl
		up.downl++
	}
}

func (h *ArrInt) sort4(down *[4]uint32) {
	v0 := h.at(down[0]).value
	v1 := h.at(down[1]).value
	v2 := h.at(down[2]).value
	v3 := h.at(down[3]).value
	if v1 < v0 {
		down[0], down[1] = down[1], down[0]
		v1, v0 = v0, v1
	}
	if v2 < v1 {
		down[1], down[2] = down[2], down[1]
		v2, v1 = v1, v2
		if v1 < v0 {
			down[0], down[1] = down[1], down[0]
			v1, v0 = v0, v1
		}
	}
	if v3 < v2 {
		down[2], down[3] = down[3], down[2]
		v3, v2 = v2, v3
		if v2 < v1 {
			down[1], down[2] = down[2], down[1]
			v2, v1 = v1, v2
			if v1 < v0 {
				down[0], down[1] = down[1], down[0]
			}
		}
	}
	h.putUnder(down[2], down[3])
	h.putUnder(down[1], down[2])
	h.putUnder(down[0], down[1])
}

func (h *ArrInt) sort3(down *[4]uint32) {
	v0 := h.at(down[0]).value
	v1 := h.at(down[1]).value
	v2 := h.at(down[2]).value
	if v1 < v0 {
		down[0], down[1] = down[1], down[0]
		v1, v0 = v0, v1
	}
	if v2 < v1 {
		down[1], down[2] = down[2], down[1]
		v2, v1 = v1, v2
		if v1 < v0 {
			down[0], down[1] = down[1], down[0]
		}
	}
	h.putUnder(down[1], down[2])
	h.putUnder(down[0], down[1])
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
	*t = intArr{ref: tm, value: tm.Value()}
	return
}

func (h *ArrInt) getChunk(i int) *[pageSize]intArr {
	return h.heap[i>>pageLog]
}
