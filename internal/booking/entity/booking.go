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
	CreatedAt     time.Time   `json:"created_at" form:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at" form:"updated_at"`
	Passengers    []Passenger `json:"passengers" form:"passengers"`
	Tickets       []Ticket    `json:"tickets" form:"tickets"`
}

type BookingCreate struct {
	Code          string    `json:"code" form:"code"`
	Total         float64   `json:"total" form:"total"`
	PaymentStatus int8      `json:"payment_status" form:"payment_status"`
	CreatedAt     time.Time `json:"created_at" form:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" form:"updated_at"`
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

func (c *BookingCreate) TableName() string {
	return "bookings"
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

// Product là một struct đại diện cho sản phẩm
type Product struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	Description string
	Price       float64
	Stock       int
	CategoryID  uint
}

// Order là một struct đại diện cho đơn hàng
type Order struct {
	ID          uint `gorm:"primaryKey"`
	CustomerID  uint
	OrderDate   string
	TotalAmount float64
	Status      string
}

// OrderItem là một struct đại diện cho mặt hàng trong đơn hàng
type OrderItem struct {
	ID        uint `gorm:"primaryKey"`
	OrderID   uint
	ProductID uint
	Quantity  int
	Price     float64
}
