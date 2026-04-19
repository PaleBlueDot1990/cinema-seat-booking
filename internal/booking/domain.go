package booking

import (
	"context"
	"errors"
	"time"
)

var (
	ErrSeatAlreadyBooked = errors.New("seat is already taken")
)

type Booking struct {
	ID        string 
	MovieID   string 
	SeatID    string 
	UserID    string 
	Status    string 
	ExpiresAt time.Time
}

type BookingStore interface {
	Book(b Booking) (Booking, error)
	ListBookings(movieID string) []Booking
	Confirm(context.Context, string, string) (Booking, error)
	Release(context.Context, string, string) error
}