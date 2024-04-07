package common

import "time"

type SQLModel struct {
	Id        int        `json:"-" gorm:"column:id;"`
	FakeId    *UID       `json:"id" gorm:"-"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"updated_at"`
}

func (m *SQLModel) GenUid(dbType int) {
	if m == nil {
		return
	}

	uid := NewUID(uint32(m.Id), dbType, 1)
	m.FakeId = &uid
}

func NewSQLModel() SQLModel {
	now := time.Now().UTC()

	return SQLModel{
		CreatedAt: &now,
		UpdatedAt: &now,
	}
}
