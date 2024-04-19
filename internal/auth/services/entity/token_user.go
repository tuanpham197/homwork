package entity

import (
	commonapp "Ronin/common"
)

type TokenUser struct {
	Id                 int    `json:"id" gorm:"column:id;"`
	UserId             string `json:"user_id" gorm:"column:user_id;"`
	Revoked            bool   `json:"revoked" gorm:"column:revoked;"`
	DeviceId           string `json:"device_id" gorm:"column:device_id;"`
	TokenSign          string `json:"token_sign" gorm:"column:token_sign"`
	commonapp.SQLModel `json:",inline"`
}

func (TokenUser) TableName() string {
	return "token_user"
}

type TokenUserCreate struct {
	UserId             string `json:"user_id" gorm:"column:user_id;"`
	Revoked            bool   `json:"revoked" gorm:"column:revoked;"`
	DeviceId           string `json:"device_id" gorm:"column:device_id;"`
	TokenSign          string `json:"token_sign" gorm:"column:token_sign"`
	commonapp.SQLModel `json:",inline"`
}

func (e *TokenUserCreate) TableName() string {
	return TokenUser{}.TableName()
}

type TokenUserUpdate struct {
	UserId    string `json:"user_id" gorm:"column:user_id;"`
	DeviceId  string `json:"device_id" gorm:"column:device_id;"`
	TokenSign string `json:"token_sign" gorm:"column:token_sign"`
}

func (e *TokenUserUpdate) TableName() string {
	return TokenUser{}.TableName()
}
