package service

type TokenService struct {
	*Service
}

func NewTokenService(s *Service) *TokenService {
	return &TokenService{
		Service: s,
	}
}

func (s *TokenService) GetTokenByBizType(bizType string) (string, error) {
	token, err := s.Repo.Token.GetTokenByBizType(s.Repo.DB.GetConn(), bizType)
	if err != nil {
		return "", err
	}
	return token, nil
}
