package entity

import (
	commonapp "Ronin/common"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Key struct {
	Id               int        `json:"id" gorm:"column:id;"`
	UserId           *uuid.UUID `json:"user_id" gorm:"column:user_id"`
	User             User       `json:"user" gorm:"foreignKey:Id;references:UserId"`
	PublicKey        string     `json:"public_key" gorm:"column:public_key;"`
	PrivateKey       string     `json:"private_key" gorm:"column:private_key"`
	RefreshToken     *string    `json:"refresh_token" gorm:"column:refresh_token;"`
	RefreshTokenUsed *[]string  `json:"refresh_token_used" gorm:"column:refresh_token_used"`
}

func (k *Key) TableName() string {
	return "keys"
}

type Permission struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Name      string     `json:"name"`
	GuardName string     `json:"guard_name"`
	Role      []Role     `json:"roles" gorm:"many2many:role_has_permissions;"`
	CreatedAt *time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Role struct {
	ID          uint         `json:"id" gorm:"primaryKey"`
	Name        string       `json:"name"`
	GuardName   string       `json:"guard_name" mapstructure:"guard_name"`
	Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:role_has_permissions;"`
	CreatedAt   *time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   *time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
}

type User struct {
	Id                 uuid.UUID  `json:"id" gorm:"size:191;primary_key"`
	UserName           string     `json:"username" gorm:"size:255;"`
	LastName           string     `json:"last_name" gorm:"size:255;"`
	FirstName          string     `json:"first_name" gorm:"size:255;"`
	Email              string     `json:"email" gorm:"size:255;"`
	Password           string     `json:"-" gorm:"size:255;"`
	Birthday           *time.Time `json:"birthday"`
	Salt               string     `json:"-" gorm:"size:30;"`
	Avatar             string     `json:"avatar" gorm:"size:255;"`
	commonapp.SQLModel `json:",inline"`
	Roles              []Role       `json:"roles,omitempty" gorm:"many2many:model_has_roles;joinForeignKey:ModelId"`
	Permissions        []Permission `json:"permissions,omitempty" gorm:"many2many:model_has_permissions;joinForeignKey:ModelId"`
}

func NewUser(id uuid.UUID, username, lastName, firstName, email, password string, result *time.Time) User {
	return User{
		Id:        id,
		UserName:  username,
		LastName:  lastName,
		FirstName: firstName,
		Email:     email,
		Password:  password,
		Birthday:  result,
	}
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.Id = uuid.New()

	return
}

func (u *User) HasRole(targetRole string) bool {
	targetRole = strings.ToLower(targetRole)

	for _, role := range u.Roles {
		if strings.ToLower(role.Name) == targetRole {
			return true
		}
	}
	return false
}

func (u *User) HasPermission(permissionTarget string) bool {
	permissionTarget = strings.ToLower(permissionTarget)

	// Has direct permission in table model_has_permission
	for _, permission := range u.Permissions {
		if strings.ToLower(permission.Name) == permissionTarget {
			return true
		}
	}

	// Has role include permission
	for _, role := range u.Roles {
		for _, permission := range role.Permissions {
			if u.HasPermission(permission.Name) {
				return true
			}
		}
	}

	return false
}

func ErrInvalidPasswordOrEmail() error {
	return errors.New("invalid password or email")
}
