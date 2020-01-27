package model

import (
	"time"
)

type User struct {
	Id        int64     `json:"id,omitempty" xorm:"pk BIGINT(20)"`
	CreatedAt time.Time `json:"created_at,omitempty" xorm:"not null created DATETIME"`
	UpdatedAt time.Time `json:"updated_at,omitempty" xorm:"not null updated DATETIME"`
	Mobile    string    `json:"mobile,omitempty" xorm:"not null default '' VARCHAR(20)"`
}
