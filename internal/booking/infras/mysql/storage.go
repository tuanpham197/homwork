package mysql

import (
	"Ronin/internal/booking/entity"
	"Ronin/internal/booking/service"
	"context"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const tableName = "bookings"

const prefix = "look_booking:"

type mysqlRepo struct {
	db  *gorm.DB
	rds *redis.Client
}

func (m mysqlRepo) Create(ctx context.Context, data *entity.BookingCreate) (*entity.Booking, error) {
	//TODO implement me
	panic("implement me")
}

func (m mysqlRepo) CreateTicket(ctx context.Context, data *entity.TicketCreate) (*entity.Ticket, error) {
	//TODO implement me
	panic("implement me")
}

func NewMySQLRepo(db *gorm.DB, rds *redis.Client) service.BookingRepository {
	return mysqlRepo{db: db, rds: rds}
}
