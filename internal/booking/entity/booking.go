package entity

import (
	"errors"
	"time"
)

type Booking struct {
	Id          int64     `json:"id" form:"id"`
	Code        string    `json:"code" form:"code"`
	Status      int8      `json:"status" form:"status"`
	CustomerId  int64     `json:"customer_id" form:"customer_id"`
	FlightId    int64     `json:"flight_id" form:"flight_id"`
	BookingDate time.Time `json:"booking_date" form:"booking_date"`
	Class       int       `json:"class" form:"class"`
	CreatedAt   time.Time `json:"created_at" form:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" form:"updated_at"`
	Customer    *Customer `json:"customer" form:"customer"`
	Tickets     []Ticket  `json:"tickets" form:"tickets"`
}

type Customer struct {
	Id      int64  `json:"id" form:"id"`
	Name    string `json:"name" form:"name"`
	Phone   string `json:"phone" form:"phone"`
	Email   string `json:"email" form:"email"`
	Gender  string `json:"gender" form:"gender"`
	Age     int    `json:"age" form:"age"`
	Country string `json:"country" form:"country"`
}

type BookingCreate struct {
	Seat       string `json:"seat" form:"seat"`
	FlightId   int64  `json:"flight_id" form:"flight_id"`
	CustomerId int64  `json:"customer_id" form:"customer_id"`
}

type Ticket struct {
	Id           int64     `json:"id" form:"id"`
	Price        int       `json:"price" form:"price"`
	Status       int64     `json:"status" form:"status"`
	SeatNumber   string    `json:"seat_number" form:"seat_number"`
	BoardingTime time.Time `json:"boarding_time" form:"boarding_time"`
	Gate         int32     `json:"gate" form:"gate"`
	Name         string    `json:"name" form:"name"`
	FlightId     int64     `json:"flight_id" form:"flight_id"`
	BookingId    int64     `json:"booking_id" form:"booking_id"`
	CustomerId   int64     `json:"customer_id" form:"customer_id"`
	Customer     *Customer `json:"customer" form:"customer"`
}

type TicketCreate struct {
	SeatNumber   string    `json:"seat_number" form:"seat_number"`
	BoardingTime time.Time `json:"boarding_time" form:"boarding_time"`
	Name         string    `json:"name" form:"name"`
	FlightId     int64     `json:"flight_id" form:"flight_id"`
	BookingId    int64     `json:"booking_id" form:"booking_id"`
	CustomerId   int64     `json:"customer_id" form:"customer_id"`
	Customer     *Customer `json:"customer" form:"customer"`
}

func (c *BookingCreate) TableName() string {
	return "bookings"
}

func (r *BookingCreate) Validate() error {

	return nil
}

var (
	ErrorTotalInvalid = errors.New("total invalid")
)
