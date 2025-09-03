package repository

import (
	"airbnb/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingRepo struct {
	DB *gorm.DB
}

func NewBookingRepo(db *gorm.DB) *BookingRepo {
	return &BookingRepo{
		DB: db,
	}
}

func (r *BookingRepo) CreateBooking(ctx context.Context, booking *models.Booking) error {
	return r.DB.WithContext(ctx).Create(booking).Error
}

func (r *BookingRepo) GetBookingByID(ctx context.Context, id uuid.UUID) (*models.Booking, error) {
	var booking models.Booking
	if err := r.DB.WithContext(ctx).First(&booking, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepo) GetBookingsByUserID(ctx context.Context, userID uuid.UUID) ([]models.Booking, error) {
	var bookings []models.Booking
	if err := r.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *BookingRepo) GetBookingsByPropertyID(ctx context.Context, propertyID uuid.UUID) ([]models.Booking, error) {
	var bookings []models.Booking
	if err := r.DB.WithContext(ctx).Where("property_id = ?", propertyID).Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *BookingRepo) CancelBooking(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Delete(&models.Booking{}, "id = ?", id).Error
}

func (r *BookingRepo) ConfirmBooking(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Model(&models.Booking{}).
		Where("id = ?", id).
		Update("status", models.Confirmed).Error
}

func (r *BookingRepo) GetUserBookings(ctx context.Context, userID uuid.UUID) ([]models.UserGetBooking, error) {
	var bookings []models.UserGetBooking
	err := r.DB.WithContext(ctx).
		Table("bookings").
		Select("bookings.id as booking_id, bookings.property_id, properties.name as property_name, bookings.status").
		Joins("JOIN properties ON bookings.property_id = properties.id").
		Where("bookings.user_id = ?", userID).
		Scan(&bookings).Error
	return bookings, err
}

func (r *BookingRepo) GetPropertyBookings(ctx context.Context, ownerID uuid.UUID) ([]models.PropertyBooking, error) {
	var bookings []models.PropertyBooking
	err := r.DB.WithContext(ctx).
		Table("bookings").
		Select("bookings.id as booking_id, bookings.property_id, bookings.user_id, bookings.status").
		Joins("JOIN properties ON bookings.property_id = properties.id").
		Where("properties.owner_id = ?", ownerID).
		Scan(&bookings).Error
	return bookings, err
}

func (r *BookingRepo) GetUserBookingByID(ctx context.Context, bookingID uuid.UUID) (*models.UserGetBooking, error) {
	var booking models.UserGetBooking
	err := r.DB.WithContext(ctx).
		Table("bookings").
		Select("bookings.id as booking_id, properties.id as property_id, properties.name as property_name, bookings.status").
		Joins("JOIN properties ON bookings.property_id = properties.id").
		Where("bookings.id = ?", bookingID).
		Scan(&booking).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepo) GetPropertyBookingByID(ctx context.Context, ID uuid.UUID) (*models.PropertyBooking, error) {
	var booking models.PropertyBooking
	err := r.DB.WithContext(ctx).
		Table("bookings").
		Select("bookings.id as booking_id, bookings.property_id, bookings.user_id, bookings.status").
		Where("bookings.id = ?", ID).
		Scan(&booking).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}
