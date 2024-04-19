package entity

import (
	commonapp "Ronin/common"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

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
	Roles              []entity.Role       `json:"roles,omitempty" gorm:"many2many:model_has_roles;joinForeignKey:ModelId"`
	Permissions        []entity.Permission `json:"permissions,omitempty" gorm:"many2many:model_has_permissions;joinForeignKey:ModelId"`
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
