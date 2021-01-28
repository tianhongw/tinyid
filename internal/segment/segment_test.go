package segment

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextID(t *testing.T) {
	tests := []struct {
		currentId int64
		maxId     int64
		loadingId int64
		delta     int
		remainder int
		want      []ID
	}{
		{
			currentId: 0,
			maxId:     10,
			loadingId: 8,
			delta:     2,
			remainder: 0,
			want: []ID{
				{Value: 2, Status: StatusOK},
				{Value: 4, Status: StatusOK},
				{Value: 6, Status: StatusOK},
				{Value: 8, Status: StatusNeedLoad},
				{Value: 10, Status: StatusNeedLoad},
				{Value: 12, Status: StatusOver},
			},
		},
		{
			currentId: 0,
			maxId:     10,
			loadingId: 8,
			delta:     2,
			remainder: 1,
			want: []ID{
				{Value: 1, Status: StatusOK},
				{Value: 3, Status: StatusOK},
				{Value: 5, Status: StatusOK},
				{Value: 7, Status: StatusOK},
				{Value: 9, Status: StatusNeedLoad},
				{Value: 11, Status: StatusOver},
			},
		},
	}

	for _, tt := range tests {
		seg := NewSegment(tt.currentId, tt.maxId, tt.loadingId,
			tt.delta, tt.remainder)
		for _, expected := range tt.want {
			got := seg.NextID()
			assert.Equal(t, expected.Value, got.Value)
			assert.Equal(t, expected.Status, got.Status)
		}
	}
}
