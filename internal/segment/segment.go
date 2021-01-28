package segment

import (
	"sync"
	"sync/atomic"
)

type Segment struct {
	currentId int64
	maxId     int64
	loadingId int64
	delta     int
	remainder int
	initOnce  sync.Once
}

func NewSegment(currentId, maxId, loadingId int64, delta, remainder int) *Segment {
	seg := &Segment{
		currentId: currentId,
		maxId:     maxId,
		loadingId: loadingId,
		delta:     delta,
		remainder: remainder,
	}

	return seg
}

func (s *Segment) NextID() *ID {
	s.init()

	v := atomic.AddInt64(&s.currentId, int64(s.delta))

	id := &ID{
		Value: v,
	}

	switch {
	case v > s.maxId:
		id.Status = StatusOver
	case v >= s.loadingId:
		id.Status = StatusNeedLoad
	default:
		id.Status = StatusOK
	}

	return id
}

func (s *Segment) Valid() bool {
	return atomic.LoadInt64(&s.currentId) <= s.maxId
}

func (s *Segment) init() {
	s.initOnce.Do(func() {
		if s.isValidCurrentId() {
			return
		}

		for i := 0; i <= s.delta; i++ {
			s.currentId += 1
			if s.isValidCurrentId() {
				s.currentId -= int64(s.delta)
				return
			}
		}
	})
}

func (s *Segment) isValidCurrentId() bool {
	return (s.currentId % int64(s.delta)) == int64(s.remainder)
}

type IDStatus int8

const (
	StatusOK       IDStatus = 0
	StatusNeedLoad IDStatus = 1
	StatusOver     IDStatus = 2
)

type ID struct {
	Value  int64
	Status IDStatus
}
