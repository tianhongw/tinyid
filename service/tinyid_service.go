package service

import (
	"fmt"
	"time"

	"github.com/tianhongw/tinyid/internal/generator"
	"github.com/tianhongw/tinyid/internal/generator/cached"
	"github.com/tianhongw/tinyid/internal/segment"
	"github.com/tianhongw/tinyid/pkg/retry"
)

const (
	defaultRetryTimes = 3
)

type TinyIdService struct {
	*Service

	idFactory *generator.GeneratorFactory
}

func NewTinyIdService(s *Service) *TinyIdService {
	tinyIdService := &TinyIdService{
		Service: s,
	}

	tinyIdService.idFactory = generator.NewGeneratorFactory(func(bizType string) (generator.Generator, error) {
		return cached.NewCachedGenerator(bizType, tinyIdService)
	})

	return tinyIdService
}

func (s *TinyIdService) GetNextSegment(bizType string) (*segment.Segment, error) {
	var seg *segment.Segment
	db := s.Service.Repo.DB.GetConn()

	segmentFunc := func() error {
		tinyId, err := s.Service.Repo.TinyId.GetTinyIdByBizType(db, bizType)
		if err != nil {
			return err
		}

		newMaxId := tinyId.MaxId + int64(tinyId.Step)
		updated, err := s.Service.Repo.TinyId.UpdateTinyId(db, tinyId.ID, newMaxId,
			tinyId.MaxId, tinyId.Version)
		if err != nil {
			return fmt.Errorf("update database failed: %v", err)
		}
		if !updated {
			return fmt.Errorf("update database failed: conflicted")
		}

		currentId := tinyId.MaxId
		loadingId := currentId + int64(tinyId.Step*s.Conf.LoadingPercent/100)

		seg = segment.NewSegment(currentId, newMaxId, loadingId,
			tinyId.Delta, tinyId.Remainder)
		return nil
	}

	err := retry.Times(defaultRetryTimes).
		Interval(100 * time.Millisecond).
		Do(segmentFunc)

	if err != nil {
		s.Service.Logger.Errorf("get next segment for bizType: %s failed after %d tries: %v",
			bizType, defaultRetryTimes, err)
	}

	return seg, err
}

func (s *TinyIdService) GetGenerator(bizType string) (generator.Generator, error) {
	return s.idFactory.GetGenerator(bizType)
}
