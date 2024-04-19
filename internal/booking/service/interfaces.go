package service

import (
	"Ronin/internal/booking/entity"
	"context"
)

type BookingUseCase interface {
	CreateBooking(ctx context.Context, data *entity.BookingCreate) (*entity.Booking, error)
	RemoveBooking(ctx context.Context, id int) (*entity.Booking, error)
	GetListBooking(ctx context.Context) ([]*entity.Booking, error)
}

type BookingRepository interface {
	Create(ctx context.Context, data *entity.BookingCreate) (entity.Booking, error)
	FindAll(ctx context.Context) ([]entity.Booking, error)
	Delete(ctx context.Context, id int) error
}
