package mysql

import (
	"Ronin/internal/booking/entity"
	"Ronin/internal/booking/service"
	"context"
	"github.com/redis/go-redis/v9"
	"time"

	"gorm.io/gorm"
)

const tableName = "bookings"

const prefix = "look_booking:"

type mysqlRepo struct {
	db  *gorm.DB
	rds *redis.Client
}

func (m mysqlRepo) FindAll(ctx context.Context) ([]entity.Booking, error) {
	//TODO implement me
	panic("implement me")
}

func (m mysqlRepo) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

func (m mysqlRepo) Create(ctx context.Context, data *entity.BookingCreate) (entity.Booking, error) {
	booking := entity.BookingCreate{
		Code:          "BK0001",
		PaymentStatus: 1,
		Total:         1000,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err := m.db.Create(&booking).Error

	if err != nil {
		return entity.Booking{}, err
	}

	return entity.Booking{
		Code:          "BK0001",
		PaymentStatus: booking.PaymentStatus,
		Total:         booking.Total,
		CreatedAt:     booking.CreatedAt,
		UpdatedAt:     booking.UpdatedAt,
	}, nil
}

func NewMySQLRepo(db *gorm.DB, rds *redis.Client) service.BookingRepository {
	return mysqlRepo{db: db, rds: rds}
}
