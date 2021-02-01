package model

import "time"

type Token struct {
	Id          uint64    `db:"id"`
	BizType     string    `db:"biz_type"`
	Token       string    `db:"token"`
	Description string    `db:"description"`
	CreateTime  time.Time `db:"create_time"`
	UpdateTime  time.Time `db:"update_time"`
}
