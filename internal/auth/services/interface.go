package service

import (
	"Ronin/internal/auth/services/entity"
	"Ronin/internal/auth/services/request"
	"context"
)

type UserUseCase interface {
	Login(ctx context.Context, userLogin request.UserLogin) (*request.UserLoginResponse, error)
	Register(ctx context.Context, userRegister request.UserRegister) (*request.UserLoginResponse, error)
	RefreshToken(ctx context.Context, token request.RefreshTokenRequest) (*request.UserLoginResponse, error)
	GetInfoUser(ctx context.Context, id string) (*entity.User, error)
	RevokeToken(ctx context.Context, req *entity.TokenUserUpdate) error
}

type UserRepository interface {
	Insert(ctx context.Context, userInfo request.UserRegister) (*entity.User, error)
	GetOne(ctx context.Context, conditions map[string]interface{}, rela ...string) (*entity.User, error)
	GetAll(ctx context.Context) ([]entity.User, error)
	Delete(ctx context.Context, id string) (bool, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	UpdateRefreshTokenUsed(ctx context.Context, refreshToken string, keyId int) (*entity.Key, error)
	InsertTokenUser(ctx context.Context, tokenUser *entity.TokenUserCreate) error
	UpdateUserToken(ctx context.Context, conditions, values map[string]interface{}) error
}
