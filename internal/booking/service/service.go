package service

import (
	"Ronin/internal/booking/entity"
	"context"
)

type service struct {
	repository BookingRepository
}

func NewService(repository BookingRepository) service {
	return service{repository: repository}
}

func (s service) CreateBooking(ctx context.Context, data *entity.BookingCreate) (*entity.Booking, error) {
	booking, err := s.repository.Create(ctx, data)

	if err != nil {
		return nil, err
	}

	return &booking, nil
}

func (s service) RemoveBooking(ctx context.Context, id int) (*entity.Booking, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) GetListBooking(ctx context.Context) ([]*entity.Booking, error) {
	//TODO implement me
	panic("implement me")
}
