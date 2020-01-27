package model

import (
	"time"
)

type User struct {
	Id        int64     `json:"id,omitempty" xorm:"pk BIGINT(20)"`
	CreatedAt time.Time `json:"created_at,omitempty" xorm:"not null created DATETIME"`
	UpdatedAt time.Time `json:"updated_at,omitempty" xorm:"not null updated DATETIME"`
	DeletedAt int64     `json:"deleted_at,omitempty" xorm:"not null default 0 BIGINT(20)"`
	Mobile    string    `json:"mobile,omitempty" xorm:"not null default '' VARCHAR(20)"`
	Passwd    string    `json:"passwd,omitempty" xorm:"not null default '' VARCHAR(128)"`
}
