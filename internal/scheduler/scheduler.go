package scheduler

import (
	"sync"
	"time"
	"github.com/ichbingautam/distributed-task-scheduler/internal/core"
	"github.com/ichbingautam/distributed-task-scheduler/internal/storage"
)

type Scheduler struct {
	store       storage.Store
	heap        *Heap
	addChan     chan *core.Task
	removeChan  chan string
	triggerChan chan struct{}
	stopChan    chan struct{}
	mu          sync.Mutex
}

func NewScheduler(store storage.Store) *Scheduler {
	return &Scheduler{
		store:       store,
		heap:        NewHeap(),
		addChan:     make(chan *core.Task, 1e6),
		removeChan:  make(chan string, 1000),
		triggerChan: make(chan struct{}, 1),
		stopChan:    make(chan struct{}),
	}
}

func (s *Scheduler) Schedule(task *core.Task) {
	s.addChan <- task
}

func (s *Scheduler) Remove(taskID string) {
	s.removeChan <- taskID
}

func (s *Scheduler) Start() <-chan *core.Task {
	out := make(chan *core.Task, 1e5)

	go func() {
		defer close(out)

		var timer *time.Timer
		for {
			select {
			case task := <-s.addChan:
				s.mu.Lock()
				s.heap.Push(task)
				s.mu.Unlock()
				s.trigger()

			case <-s.removeChan:
				s.mu.Lock()
				// Implementation of removal logic
				s.mu.Unlock()

			case <-s.triggerChan:
				s.mu.Lock()
				next := s.heap.Peek()
				s.mu.Unlock()

				if next == nil {
					continue
				}

				if timer != nil {
					timer.Stop()
				}
				delay := time.Until(next.ScheduledAt)
				timer = time.NewTimer(delay)

				select {
				case <-timer.C:
					s.mu.Lock()
					task := s.heap.Pop()
					s.mu.Unlock()
					out <- task
				case <-s.stopChan:
					if timer != nil {
						timer.Stop()
					}
					return
				}

			case <-s.stopChan:
				if timer != nil {
					timer.Stop()
				}
				return
			}
		}
	}()

	return out
}

func (s *Scheduler) trigger() {
	select {
	case s.triggerChan <- struct{}{}:
	default:
	}
}

func (s *Scheduler) Stop() {
	close(s.stopChan)
}