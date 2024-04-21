package mysql

import (
	"Ronin/internal/auth/services"
	"Ronin/internal/auth/services/entity"
	"Ronin/internal/auth/services/request"
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type mysqlRepo struct {
	db   *gorm.DB
	reds *redis.Client
}

func (repo mysqlRepo) UpdateUserToken(ctx context.Context, conditions, values map[string]interface{}) error {
	return repo.db.Debug().Table(entity.TokenUser{}.TableName()).Where(conditions).Updates(values).Error
}

func (repo mysqlRepo) GetTokenUser(ctx context.Context, conditions map[string]interface{}) (*entity.TokenUser, error) {
	var tokenUser entity.TokenUser
	if err := repo.db.Where(conditions).First(&tokenUser).Error; err != nil {
		return nil, err
	}
	return &tokenUser, nil
}

func (repo mysqlRepo) InsertTokenUser(ctx context.Context, tokenUser *entity.TokenUserCreate) error {

	if err := repo.db.Create(tokenUser).Error; err != nil {
		return err
	}

	return nil
}

func NewMySQLRepo(db *gorm.DB, reds *redis.Client) service.UserRepository {
	return mysqlRepo{db: db, reds: reds}
}

func (repo mysqlRepo) GetOne(ctx context.Context, conditions map[string]interface{}, relations ...string) (*entity.User, error) {
	var user entity.User
	db := repo.db
	if len(relations) > 0 {
		for _, pre := range relations {
			db.Preload(pre)
		}
	}
	err := db.Where(conditions).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo mysqlRepo) GetAll(ctx context.Context) ([]entity.User, error) {
	panic("implement me")
}

func (repo mysqlRepo) Insert(ctx context.Context, userRegister request.UserRegister) (*entity.User, error) {

	user := entity.User{
		UserName:  userRegister.UserName,
		LastName:  userRegister.LastName,
		FirstName: userRegister.FirstName,
		Email:     userRegister.Email,
		Password:  userRegister.Password,
		Salt:      userRegister.Salt,
	}

	result := repo.db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
func (repo mysqlRepo) Delete(ctx context.Context, id string) (bool, error) {
	panic("implement me")
}

func (repo mysqlRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := repo.db.Debug().
		Preload("Roles", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo mysqlRepo) AssignRoleUser(ctx context.Context, reqRole request.AssignRoleUser, userId string) (bool, error) {
	db := repo.db
	user, errUser := repo.GetOne(ctx, map[string]interface{}{
		"id": userId,
	}, "Roles", "Permissions")
	if errUser != nil || user == nil {
		return false, errUser
	}

	var roles []entity.Role
	errRole := db.Find(&roles, reqRole.Roles).Error
	if errRole != nil || len(roles) < 1 {
		return false, errors.New("not found role to assign")
	}

	errAppend := db.Model(user).Omit("Roles.*").Association("Roles").Replace(&roles)
	if errAppend != nil {
		return false, errAppend
	}

	return true, nil
}

func (repo mysqlRepo) StoreKey(ctx context.Context, keyRequest request.KeyRequest) (*entity.Key, error) {
	key := entity.Key{
		PrivateKey:       keyRequest.PrivateKey,
		PublicKey:        keyRequest.PublicKey,
		UserId:           keyRequest.UserId,
		RefreshToken:     keyRequest.RefreshToken,
		RefreshTokenUsed: keyRequest.RefreshTokenUsed,
	}

	if err := repo.db.Create(&key).Error; err != nil {
		return nil, err
	}
	return &key, nil
}

func (repo mysqlRepo) UpdateRefreshTokenUsed(ctx context.Context, refreshToken string, keyId int) (*entity.Key, error) {
	var key entity.Key
	errGetKey := repo.db.Where("id = ?", keyId).First(&key).Error
	if errGetKey != nil {
		return nil, errGetKey
	}
	refreshTokenUsed := *key.RefreshTokenUsed
	refreshTokenUsed = append(refreshTokenUsed, refreshToken)
	key.RefreshTokenUsed = &refreshTokenUsed
	errSave := repo.db.Save(&key).Error
	if errSave != nil {
		return nil, errSave
	}
	return &key, nil
}
