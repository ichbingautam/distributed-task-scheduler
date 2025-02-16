package scheduler

import (
	"container/heap"
	"github.com/ichbingautam/distributed-task-scheduler/internal/core"
)

// type core.Task struct {
// 	ID          string
// 	Execute     func() error
// 	ScheduledAt time.Time
// 	RetryPolicy core.RetryPolicy
// 	Attempts    int
// }

type taskHeap []*core.Task

func (h taskHeap) Len() int            { return len(h) }
func (h taskHeap) Less(i, j int) bool  { return h[i].ScheduledAt.Before(h[j].ScheduledAt) }
func (h taskHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *taskHeap) Push(x interface{}) { *h = append(*h, x.(*core.Task)) }
func (h *taskHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type Heap struct {
	heap *taskHeap
}

func NewHeap() *Heap {
	h := &taskHeap{}
	heap.Init(h)
	return &Heap{heap: h}
}

func (h *Heap) Push(task *core.Task) {
	heap.Push(h.heap, task)
}

func (h *Heap) Pop() *core.Task {
	return heap.Pop(h.heap).(*core.Task)
}

func (h *Heap) Peek() *core.Task {
	if len(*h.heap) == 0 {
		return nil
	}
	return (*h.heap)[0]
}