/*
Package fastheap provides a heap implementation specialized on numeric keys and arbitrary stored values.

It is faster than containers/heap cause it doesn't use interface call for comparisons, and it manipulates storage by itself.

Also, it uses 4 way heap instead of binary heap. May be it makes a deal with CPU cache :)
*/
package fastheap

import (
	"errors"
)

// ErrInsert - error returned by heap.Insert(item) if item index is not zero
var ErrInsert = errors.New("Could not insert heap value with index != 0")
// ErrRemove - error returned by heap.Remove(item) and heap.Reset(item) if item index is out of range
var ErrRemove = errors.New("Could not remove heap value with wrong index")
// ErrPop - error returned by heap.Pop(item) if heap is empty
var ErrPop = errors.New("Could not pop from empty heap")

const (
	pageLog = 8
	pageSize = 1 << pageLog
	pageMask = pageSize - 1
)
