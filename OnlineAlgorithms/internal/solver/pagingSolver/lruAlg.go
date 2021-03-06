package pagingsolver

import (
	"OnlineAlgorithms/internal/solver"
	"container/heap"
	"fmt"
)

type LRUMem struct {
	mem     int
	lastReq int
	index   int
}

type PriorityQueue []*LRUMem

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].lastReq > pq[j].lastReq
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*LRUMem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) update(item *LRUMem) {
	item.lastReq++
}

type LRUAlg struct {
	memory PriorityQueue
	size   int
	debug  bool
}

func LRUAlg_Create(size int, debug bool) *LRUAlg {
	lru := &LRUAlg{size: size, memory: make(PriorityQueue, 0), debug: debug}
	heap.Init(&lru.memory)

	return lru
}

func (alg *LRUAlg) UpdateMemory(request int) bool {
	isFound := alg.find(request)
	solver.DebugPrint(fmt.Sprint(alg.unpackMemory()), alg.debug)
	heap.Init(&alg.memory)
	if !isFound {
		solver.DebugPrint(fmt.Sprint(" ## FAULT "), alg.debug)
		solver.DebugPrint(fmt.Sprint(" HAVE TO INSERT ", request, " ## "), alg.debug)
		if alg.memory.Len() >= alg.size {
			x := heap.Pop(&alg.memory).(*LRUMem)
			solver.DebugPrint(fmt.Sprint(" ## POPPING ", x.mem, " ## "), alg.debug)
		}
		heap.Push(&alg.memory, &LRUMem{mem: request, lastReq: 0})
		solver.DebugPrint(fmt.Sprint(" =>> ", alg.unpackMemory()), alg.debug)
	} else {
		solver.DebugPrint(fmt.Sprint(" ## FOUND ", request, " REQUEST SERVED ## =>> ", alg.unpackMemory()), alg.debug)
	}
	heap.Init(&alg.memory)
	solver.DebugPrint(fmt.Sprintln(), alg.debug)
	return isFound
}

func (alg *LRUAlg) find(request int) bool {
	ret := false
	for _, n := range alg.memory {
		if n.mem == request {
			ret = true
			continue
		}
		alg.memory.update(n)
	}
	return ret
}

func (alg *LRUAlg) unpackMemory() [][2]int {
	mem := make([][2]int, 0)

	for _, n := range alg.memory {
		mem = append(mem, [2]int{n.mem, n.lastReq})
	}

	return mem
}
