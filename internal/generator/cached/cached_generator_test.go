package cached

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tianhongw/tinyid/internal/segment"
)

const (
	testBizType = "test"
)

func TestNextID(t *testing.T) {
	cg, err := NewCachedGenerator(testBizType, &fakeSegmentService{})
	assert.Nil(t, err)
	expected := []int64{2, 4, 6, 8, 10}
	for _, id := range expected {
		got, err := cg.NextID()
		assert.Nil(t, err)
		assert.Equal(t, id, got)
	}
}

func TestNextBatchIDs(t *testing.T) {
	cg, err := NewCachedGenerator(testBizType, &fakeSegmentService{})
	assert.Nil(t, err)
	expected := []int64{2, 4, 6, 8, 10}
	got, err := cg.NextBatchIDs(len(expected))
	assert.Nil(t, err)
	assert.Equal(t, expected, got)
}

type fakeSegmentService struct{}

func (s *fakeSegmentService) GetNextSegment(bizType string) (*segment.Segment, error) {
	return segment.NewSegment(0, 100, 80, 2, 0), nil
}
