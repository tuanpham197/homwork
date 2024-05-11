package service

import (
	"Ronin/internal/booking/entity"
	"context"
	"github.com/redis/go-redis/v9"
)

type service struct {
	repository BookingRepository
	rds        *redis.Client
}

func NewService(repository BookingRepository, rds *redis.Client) service {
	return service{repository: repository, rds: rds}
}

func (s service) CreateBooking(ctx context.Context, data *entity.BookingCreate) (*entity.Booking, error) {
	// lock by flight and seat

	// Create booking
	_, err := s.repository.Create(ctx, data)

	// Create ticket

	if err != nil {
		return nil, err
	}

	return nil, nil
}
