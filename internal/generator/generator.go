package generator

import (
	"sync"

	"github.com/tianhongw/tinyid/internal/segment"
)

type GeneratorCreateFunc func(bizType string) (Generator, error)

type Generator interface {
	NextID() (id int64, err error)
	NextBatchIDs(size int) (ids []int64, err error)
	GetCurrentSegment(bizType string) (current *segment.Segment)
	GetNextSegment(bizType string) (next *segment.Segment)
}

type GeneratorFactory struct {
	GeneratorCreator GeneratorCreateFunc
	generators       map[string]Generator
	mu               sync.Mutex
}

func NewGeneratorFactory(creator GeneratorCreateFunc) *GeneratorFactory {
	return &GeneratorFactory{
		GeneratorCreator: creator,
		generators:       make(map[string]Generator),
	}
}

func (gf *GeneratorFactory) GetGenerator(bizType string) (Generator, error) {
	gf.mu.Lock()
	defer gf.mu.Unlock()

	if generator, ok := gf.generators[bizType]; ok {
		return generator, nil
	}

	generator, err := gf.GeneratorCreator(bizType)
	if err != nil {
		return nil, err
	}

	gf.generators[bizType] = generator

	return generator, nil
}
