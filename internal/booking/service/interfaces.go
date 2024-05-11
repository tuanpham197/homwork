package service

import (
	"Ronin/internal/booking/entity"
	"context"
)

type BookingUseCase interface {
	CreateBooking(ctx context.Context, data *entity.BookingCreate) (*entity.Booking, error)
}

type BookingRepository interface {
	Create(ctx context.Context, data *entity.BookingCreate) (*entity.Booking, error)
	CreateTicket(ctx context.Context, data *entity.TicketCreate) (*entity.Ticket, error)
}
