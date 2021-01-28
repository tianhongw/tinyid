package cached

import (
	"log"
	"sync"

	"github.com/tianhongw/tinyid/internal/generator"
	"github.com/tianhongw/tinyid/internal/segment"
)

var _ generator.Generator = (*CachedGenerator)(nil)

type SegmentService interface {
	GetNextSegment(bizType string) (*segment.Segment, error)
}

type CachedGenerator struct {
	bizType        string
	segmentService SegmentService
	current        *segment.Segment
	currentMu      sync.Mutex
	next           *segment.Segment
	nextMu         sync.Mutex
	isLoadingNext  bool
}

func NewCachedGenerator(bizType string, segmentService SegmentService) (*CachedGenerator, error) {
	cg := &CachedGenerator{
		bizType:        bizType,
		segmentService: segmentService,
	}

	if err := cg.loadCurrent(); err != nil {
		return nil, err
	}

	return cg, nil
}

func (g *CachedGenerator) NextID() (int64, error) {
	for {
		if g.current == nil {
			if err := g.loadCurrent(); err != nil {
				return 0, err
			}
			continue
		}
		id := g.current.NextID()
		switch id.Status {
		case segment.StatusOK:
			return id.Value, nil
		case segment.StatusNeedLoad:
			g.loadNext()
			return id.Value, nil
		case segment.StatusOver:
			if err := g.loadCurrent(); err != nil {
				return 0, err
			}
		}
	}
}

func (g *CachedGenerator) NextBatchIDs(size int) ([]int64, error) {
	ids := make([]int64, 0, size)
	for i := 0; i < size; i++ {
		id, err := g.NextID()
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (g *CachedGenerator) GetCurrentSegment(bizType string) *segment.Segment {
	return g.current
}

func (g *CachedGenerator) GetNextSegment(bizType string) *segment.Segment {
	return g.next
}

func (g *CachedGenerator) loadCurrent() error {
	g.currentMu.Lock()
	defer g.currentMu.Unlock()

	if g.current == nil || !g.current.Valid() {
		if g.next != nil {
			g.current = g.next
			g.next = nil
			return nil
		}

		seg, err := g.segmentService.GetNextSegment(g.bizType)
		if err != nil {
			return err
		}
		g.current = seg
	}

	return nil
}

func (g *CachedGenerator) loadNext() {
	g.nextMu.Lock()
	defer g.nextMu.Unlock()

	if g.next == nil && !g.isLoadingNext {
		g.isLoadingNext = true
		go func() {
			seg, err := g.segmentService.GetNextSegment(g.bizType)
			if err != nil {
				log.Printf("load next segment failed: %v", err)
			} else {
				g.next = seg
				g.isLoadingNext = false
			}
		}()
	}
}
