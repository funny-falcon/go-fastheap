/*
fastheap is an heap implementation specialized on numeric keys and arbitrary stored values.
It is faster than containers/heap cause it doesn't use interface call for comparisons, and
it manipulates storage by itself.
Also, it uses 4 way heap instead of binary heap. May be it makes a deal with CPU cache :)
*/
package fastheap

import (
	"errors"
)

var InsertError = errors.New("Could not insert heap value with index != 0")
var RemoveError = errors.New("Could not remove heap value with wrong index")
var PopError = errors.New("Could not pop from empty heap")
