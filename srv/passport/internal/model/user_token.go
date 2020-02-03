package model

import (
	"time"
)

type UserToken struct {
	Id           int64     `json:"id,omitempty" xorm:"pk BIGINT(20)"`
	CreatedAt    time.Time `json:"created_at,omitempty" xorm:"not null created DATETIME"`
	UpdatedAt    time.Time `json:"updated_at,omitempty" xorm:"not null updated DATETIME"`
	AppId        int       `json:"app_id,omitempty" xorm:"not null default 0 INT(11)"`
	Uid          int64     `json:"uid,omitempty" xorm:"not null default 0 BIGINT(20)"`
	ExpiresIn    int64     `json:"expires_in,omitempty" xorm:"not null default 0 BIGINT(20)"`
	AccessToken  string    `json:"access_token,omitempty" xorm:"not null default '' VARCHAR(64)"`
	RefreshToken string    `json:"refresh_token,omitempty" xorm:"not null default '' VARCHAR(64)"`
	DeviceId     string    `json:"device_id,omitempty" xorm:"not null default '' VARCHAR(64)"`
}
