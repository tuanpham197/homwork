package service

import (
	"Ronin/internal/booking/entity"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type service struct {
	repository BookingRepository
	rds        *redis.Client
}

func (s service) LookProduct(ctx context.Context, id int) (*entity.Product, error) {
	key := fmt.Sprintf("product:%v", id)
	for {
		result, errLock := s.rds.SetNX(ctx, key, id, time.Second*10).Result()
		if errLock != nil || !result {
			time.Sleep(1 * time.Second)
			continue
		}

		// handle create order
		s.repository.Create(ctx, &entity.BookingCreate{})

		// release lock
		_, _ = s.rds.Expire(ctx, key, 0).Result()
		break
	}

	return nil, nil
}

func NewService(repository BookingRepository, rds *redis.Client) service {
	return service{repository: repository, rds: rds}
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
