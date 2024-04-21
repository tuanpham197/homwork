package appctx

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AppCtx interface {
	GetMainDBConnection() *gorm.DB
	SecretKey() string
	GetRedisClient() *redis.Client
}

type appCtx struct {
	db        *gorm.DB
	secretKey string
	rds       *redis.Client
}

func NewAppContext(db *gorm.DB, secretKey string, rds *redis.Client) *appCtx {
	return &appCtx{
		db:        db,
		secretKey: secretKey,
		rds:       rds,
	}
}

func (a *appCtx) GetMainDBConnection() *gorm.DB { return a.db }
func (a *appCtx) SecretKey() string             { return a.secretKey }
func (a *appCtx) GetRedisClient() *redis.Client { return a.rds }
