package model

import "time"

type TinyId struct {
	ID         uint64    `db:"id"`
	BizType    string    `db:"biz_type"`
	MaxId      int64     `db:"max_id"`
	Step       int       `db:"step"`
	Delta      int       `db:"delta"`
	Remainder  int       `db:"remainder"`
	Version    int64     `db:"version"`
	CreateTime time.Time `db:"create_time"`
	UpdateTime time.Time `db:"update_time"`
}
