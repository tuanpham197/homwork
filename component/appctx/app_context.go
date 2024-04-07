package appctx

import (
	"gorm.io/gorm"
)

type AppCtx interface {
	GetMainDBConnection() *gorm.DB
	SecretKey() string
}

type appCtx struct {
	db        *gorm.DB
	secretKey string
}

func NewAppContext(db *gorm.DB, secretKey string) *appCtx {
	return &appCtx{
		db:        db,
		secretKey: secretKey,
	}
}

func (a *appCtx) GetMainDBConnection() *gorm.DB { return a.db }
func (a *appCtx) SecretKey() string             { return a.secretKey }
