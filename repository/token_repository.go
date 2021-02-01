package repository

import (
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/tianhongw/tinyid/model"
)

type TokenRepository struct {
	*Repository
}

func NewTokenRepository(r *Repository) *TokenRepository {
	return &TokenRepository{
		Repository: r,
	}
}

func (r *TokenRepository) TableName() string {
	return "tiny_id_token"
}

func (r *TokenRepository) GetTokenByBizType(db *sqlx.DB, bizType string) (string, error) {
	sql, args := squirrel.Select("*").
		From(r.TableName()).
		Where(squirrel.Eq{
			"biz_type": bizType,
		}).MustSql()

	t := new(model.Token)

	if err := db.Get(t, sql, args...); err != nil {
		return "", err
	}

	return t.Token, nil
}
