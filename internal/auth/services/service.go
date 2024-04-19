package service

import (
	commonapp "Ronin/common"
	"Ronin/internal/auth/services/entity"
	"Ronin/internal/auth/services/request"
	"Ronin/pkg/contants"
	authUtil "Ronin/pkg/utils/auth"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strings"

	"go.uber.org/zap"
)

type Hasher interface {
	RandomStr(length int) (string, error)
	HashPassword(salt, password string) (string, error)
	CompareHashPassword(hashedPassword, salt, password string) bool
}

type service struct {
	repository UserRepository
	hasher     Hasher
	log        *zap.SugaredLogger
	rds        *redis.Client
}

func NewService(
	repository UserRepository,
	hasher Hasher,
	logger *zap.SugaredLogger,
	rds *redis.Client,
) service {
	return service{
		repository,
		hasher,
		logger,
		rds,
	}
}

func (s service) RevokeToken(ctx context.Context, req *entity.TokenUserUpdate) error {
	return s.repository.UpdateUserToken(ctx, map[string]interface{}{
		"user_id":    req.UserId,
		"token_sign": req.TokenSign,
		"device_id":  req.DeviceId,
	}, map[string]interface{}{
		"revoked": false,
	})
}

func (s service) Login(ctx context.Context, userLogin request.UserLogin) (*request.UserLoginResponse, error) {

	user, err := s.repository.GetByEmail(ctx, userLogin.Email)

	if err != nil {
		return nil, err
	}

	// TODO: Slow password compare by bcrypt
	resultHash := s.hasher.CompareHashPassword(user.Password, user.Salt, userLogin.Password)
	if !resultHash {
		return nil, entity.ErrInvalidPasswordOrEmail()
	}

	// generate access token
	// Assuming the login is successful, generate the tokens2
	payload := authUtil.Payload{
		UserID: user.Id,
		Roles:  user.Roles,
	}
	accessToken, err := authUtil.GenerateAccessToken(&payload)
	if err != nil {
		return nil, err
	}
	refreshToken, err := authUtil.GenerateRefreshToken(&payload)
	if err != nil {
		return nil, err
	}

	// save to token user
	//go func() {
	//	defer commonapp.AppRecover()
	//	_ = s.repository.InsertTokenUser(ctx, &entity.TokenUserCreate{
	//		UserId:    user.Id.String(),
	//		DeviceId:  "",
	//		Revoked:   false,
	//		TokenSign: strings.Split(accessToken, ".")[2],
	//	})
	//}()
	// set redis
	key := fmt.Sprintf("user_info:%s", user.Id)
	errRds := commonapp.SetDataToRedis(ctx, user, key, commonapp.AccessTokenExpireDuration, s.rds)
	if errRds != nil {
		return nil, errRds
	}

	// Return the tokens in the response
	return &request.UserLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserInfo: request.UserInfo{
			Id:        user.Id,
			UserName:  user.UserName,
			LastName:  user.LastName,
			Email:     user.Email,
			FirstName: user.FirstName,
			Role:      user.Roles,
		},
	}, nil
}

func (s service) Register(ctx context.Context, userRegister request.UserRegister) (*request.UserLoginResponse, error) {

	// check email exists
	user, _ := s.repository.GetByEmail(ctx, userRegister.Email)
	if user != nil {
		return nil, request.UserExistsError{}
	}

	//generate salt
	salt, errSalt := s.hasher.RandomStr(5)
	if errSalt != nil {
		return nil, constants.ErrorHandlePassword
	}

	// hash pass after call repo
	hashPass, errHash := s.hasher.HashPassword(salt, userRegister.Password)
	if errHash != nil {
		return nil, constants.ErrorHandlePassword
	}

	// insert to db
	userRegister.Password = hashPass
	userRegister.Salt = salt
	result, err := s.repository.Insert(ctx, userRegister)
	if err != nil {
		return nil, constants.ErrorDB
	}

	// Assuming the login is successful, generate the tokens
	payload := authUtil.Payload{
		UserID: result.Id,
	}
	accessToken, err := authUtil.GenerateAccessToken(&payload)
	if err != nil {
		return nil, nil
	}

	refreshToken, err := authUtil.GenerateRefreshToken(&payload)
	if err != nil {
		return nil, nil
	}

	go func() {
		defer commonapp.AppRecover()
		_ = s.repository.InsertTokenUser(ctx, &entity.TokenUserCreate{
			UserId:    result.Id.String(),
			DeviceId:  "",
			Revoked:   false,
			TokenSign: strings.Split(accessToken, ".")[2],
		})
	}()

	return &request.UserLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserInfo: request.UserInfo{
			Id:        result.Id,
			UserName:  result.UserName,
			LastName:  result.LastName,
			Email:     result.Email,
			FirstName: result.FirstName,
		},
	}, nil
}

func (s service) RefreshToken(ctx context.Context, token request.RefreshTokenRequest) (*request.UserLoginResponse, error) {
	// Verify the refresh token and extract the user ID
	claims, err := authUtil.VerifyRefreshToken(token.RefreshToken)
	if err != nil {
		return nil, request.ResponseMessageError{
			Message: "Failed to verify token",
		}
	}

	// Generate a new access token
	payload := authUtil.Payload{
		UserID: claims.UserID,
		ShopID: claims.ShopID,
	}
	accessToken, err := authUtil.GenerateAccessToken(&payload)
	if err != nil {
		return nil, request.ResponseMessageError{
			Message: "Failed to generate access token",
		}
	}
	// Generate new refresh token
	refreshToken, err := authUtil.GenerateRefreshToken(&payload)
	if err != nil {
		return nil, request.ResponseMessageError{
			Message: "Failed to generate refresh token",
		}
	}

	userInfo, err := s.repository.GetOne(ctx, map[string]interface{}{
		"id": claims.UserID.String(),
	})
	if err != nil {
		return nil, request.ResponseMessageError{
			Message: "Fail get info user",
		}
	}

	//TODO: handle revoke old acccess token

	// Return the new access token in the response
	return &request.UserLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserInfo: request.UserInfo{
			Id:        userInfo.Id,
			UserName:  userInfo.UserName,
			LastName:  userInfo.LastName,
			Email:     userInfo.Email,
			FirstName: userInfo.FirstName,
		},
	}, nil
}

func (s service) GetInfoUser(ctx context.Context, id string) (*entity.User, error) {
	user, err := s.repository.GetOne(ctx, map[string]interface{}{
		"id": id,
	}, "Roles", "Permissions")
	if err != nil {
		return nil, err
	}

	return user, nil
}
