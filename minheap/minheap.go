package minheap

import "errors"

type Pair struct {
	Key uint32
	Value uint32
}

type MinHeap struct {
	datas []Pair
	Size int
}

func New() *MinHeap {
	minHeap := new(MinHeap)
	minHeap.Size = 0
	minHeap.datas = []Pair{}
	return minHeap
}

func (h *MinHeap) Push(e Pair) {
	h.Size++
	h.shiftup(h.Size,e)
}

func (h *MinHeap) shiftup(s int, e Pair) {
	for s > 1 {
		parent := (s - 1) / 2
		if h.datas[parent].Key < e.Key {
			break
		}
		if s == h.Size {
			h.datas = append(h.datas,h.datas[parent])
		} else {
			h.datas[s] = h.datas[parent]
		}
		s = parent
	}
	if s == h.Size {
		h.datas = append(h.datas,e)
	} else {
		h.datas[s] = e
	}
}

func (h *MinHeap) Pop() (e Pair, err error) {
	if h.Size <= 0 {
		return Pair{Key:0,Value:0},errors.New("Empty minheap!")
	}
	ret := h.datas[0]
	h.Size--
	h.shiftdown(0,h.datas[h.Size])
	return ret,nil
}

func (h *MinHeap) shiftdown(i int, e Pair) {
	half := h.Size / 2
	for i < half {
		child := 2 * i + 1
		right := child + 1
		if right < h.Size && h.datas[child].Key > h.datas[right].Key {
			child = right
		}
		if e.Key < h.datas[child].Key {
			break
		}
		h.datas[i] = h.datas[child]
		i = child
	}
	h.datas[i] = e
}