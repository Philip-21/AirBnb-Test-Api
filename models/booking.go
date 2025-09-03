package models

import "github.com/google/uuid"

type UpdateBooking struct {
	CheckIn  string
	CheckOut string
	Status   string
}

type UserGetBooking struct {
	BookingID    uuid.UUID `json:"booking_id"`
	PropertyID   uuid.UUID `json:"property_id"`
	PropertyName string    `json:"property_name"`
	Status       string    `json:"status"`
}

type GetUserBookings struct {
	Bookings []UserGetBooking `json:"bookings"`
}

type PropertyBooking struct {
	BookingID  uuid.UUID `json:"booking_id"`
	PropertyID uuid.UUID `json:"property_id"`
	UserID     uuid.UUID `json:"user_id"`
	Status     string    `json:"status"`
}

type GetPropertyBookings struct {
	Bookings []PropertyBooking `json:"bookings"`
}
