package service

import "github.com/tianhongw/tinyid/internal/segment"

type TinyIdService struct {
	*Service
}

func NewTinyIdService(s *Service) *TinyIdService {
	return &TinyIdService{
		Service: s,
	}
}

func (s *TinyIdService) GetNextSegment(bizType string) (*segment.Segment, error) {
	_, err := s.Service.Repo.TinyId.GetTinyIdByBizType(bizType)
	if err != nil {
		s.Service.Logger.Errorf("get tiny id for bizType: %s failed: %v", bizType, err)
		return nil, err
	}
	return nil, nil
}
