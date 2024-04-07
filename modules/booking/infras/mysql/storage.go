package mysql

import (
	"Ronin/modules/booking/entity"
	"Ronin/modules/booking/service"
	"context"

	"gorm.io/gorm"
)

const tableName = "bookings"

type mysqlRepo struct {
	db *gorm.DB
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
	//TODO implement me
	panic("implement me")
}

func NewMySQLRepo(db *gorm.DB) service.BookingRepository {
	return mysqlRepo{db: db}
}
