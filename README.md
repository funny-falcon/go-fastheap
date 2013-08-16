# fastheap

[Documentation online](http://godoc.org/github.com/funny-falcon/go-fastheap)

**fastheap** provides fast heap with numerical keys.

## Install (with GOPATH set on your machine)
----------

```
go get github.com/funny-falcon/go-fastheap
```

##Usage
----------
```
package main

import (
  "fmt"
  "github.com/funny-falcon/go-fastheap"
)

type MyItem struct {
    myvalue string
    mykey   int64
    heapIndex int
}

type heapPointer struct {
   *MyItem
}

func (h heapPointer) Value() int64 {
    return h.mykey
}

func (h heapPointer) Index() int {
    return h.heapIndex
}

func (h heapPointer) SetIndex(i int) {
    h.heapIndex = i
}

func main() {
  heap := fastheap.IntHeap{}

  val1 := MyItem{myvalue: "hi", mykey: 3}
  val2 := MyItem{myvalue: "ho", mykey: 1}
  val3 := MyItem{myvalue: "hu", mykey: 2}

  heap.Insert(heapPointer{&val1})
  heap.Insert(heapPointer{&val2})
  heap.Insert(heapPointer{&val3})

  p, v, _ := heap.Top()
  fmt.Println("key:", v, "val:", p.(heapPointer).myvalue)
  p1, _ := heap.Pop()
  fmt.Println("key:", v, "val:", p1.(heapPointer).myvalue)
  p2, _ := heap.Pop()
  fmt.Println("key:", v, "val:", p2.(heapPointer).myvalue)
  p3, _ := heap.Pop()
  fmt.Println("key:", v, "val:", p3.(heapPointer).myvalue)

  _, err := heap.Pop()
  fmt.Println(err)
}
```

##License
----------
go-fastheap is BSD licensed.
