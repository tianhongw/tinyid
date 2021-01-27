package repository

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/tianhongw/tinyid/model"
)

type TinyIdRepository struct {
	*Repository
}

func NewTinyIdRepository(r *Repository) *TinyIdRepository {
	return &TinyIdRepository{
		Repository: r,
	}
}

func (r *TinyIdRepository) TableName() string {
	return "tiny_id"
}

func (r *TinyIdRepository) GetTinyIdByBizType(bizType string) (*model.TinyId, error) {
	sql, args := squirrel.Select("*").
		From(r.TableName()).
		Where(squirrel.Eq{
			"biz_type": bizType,
		}).MustSql()

	tinyId := new(model.TinyId)

	db := r.Repository.DB.GetConn()

	if err := db.Get(tinyId, db.Rebind(sql), args...); err != nil {
		return nil, err
	}

	return tinyId, nil
}

func (r *TinyIdRepository) UpdateTinyId(id uint64, newMaxId, oldMaxId, version int64) (bool, error) {
	sql, args, err := squirrel.Update(r.TableName()).
		Set("max_id", newMaxId).
		Set("version", version+1).
		Set("update_time", time.Now().UTC()).
		Where(squirrel.Eq{
			"id":      id,
			"max_id":  oldMaxId,
			"version": version,
		}).ToSql()

	if err != nil {
		return false, err
	}

	db := r.Repository.DB.GetConn()

	execResult, err := db.Exec(db.Rebind(sql), args...)
	if err != nil {
		return false, err
	}

	rowsAffected, err := execResult.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}
