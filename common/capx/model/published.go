package model

import (
	"time"
)

type Published struct {
	Id        int64     `json:"id,omitempty" xorm:"pk BIGINT(20)"`
	Topic     string    `json:"topic,omitempty" xorm:"not null default '' VARCHAR(256)"`
	Name      string    `json:"name,omitempty" xorm:"not null default '' VARCHAR(256)"`
	Version   int64     `json:"version,omitempty" xorm:"not null default 0 BIGINT(20)"`
	Msg       []byte    `json:"msg,omitempty" xorm:"not null VARBINARY(8192)"`
	Retries   int       `json:"retries,omitempty" xorm:"not null default 0 INT(11)"`
	CreatedAt time.Time `json:"created_at,omitempty" xorm:"not null created DATETIME"`
	UpdatedAt time.Time `json:"updated_at,omitempty" xorm:"not null updated DATETIME"`
	Status    int       `json:"status,omitempty" xorm:"not null default 0 TINYINT(4)"`
}
