package entity

import (
	"errors"
	"time"
)

type Booking struct {
	Id            int64       `json:"id" form:"id"`
	Code          string      `json:"code" form:"code"`
	Total         float64     `json:"total" form:"total"`
	PaymentStatus int8        `json:"payment_status" form:"payment_status"`
	Passengers    []Passenger `json:"passengers" form:"passengers"`
	Tickets       []Ticket    `json:"tickets" form:"tickets"`
}

type BookingCreate struct {
	Code          string  `json:"code" form:"code"`
	Total         float64 `json:"total" form:"total"`
	PaymentStatus int8    `json:"payment_status" form:"payment_status"`
}

type Passenger struct {
	Id      int64  `json:"id" form:"id"`
	Name    string `json:"name" form:"name"`
	Phone   string `json:"phone" form:"phone"`
	Email   string `json:"email" form:"email"`
	Gender  string `json:"gender" form:"gender"`
	Age     int    `json:"age" form:"age"`
	Country string `json:"country" form:"country"`
}

type Ticket struct {
	Id           int64     `json:"id" form:"id"`
	Code         string    `json:"code" form:"code"`
	Price        int       `json:"price" form:"price"`
	Status       int64     `json:"status" form:"status"`
	Seat         string    `json:"seat" form:"seat"`
	BoardingTime time.Time `json:"boarding_time" form:"boardingTime"`
	Gate         int32     `json:"gate" form:"gate"`
	From         string    `json:"from" form:"from"`
	To           string    `json:"to" form:"to"`
	Name         string    `json:"name" form:"name"`
	FlightId     int64     `json:"flightId" form:"flightId"`
}

func (r *BookingCreate) Validate() error {
	if r.Total < 0 {
		return ErrorTotalInvalid
	}

	return nil
}

var (
	ErrorTotalInvalid = errors.New("Total invalid")
)
