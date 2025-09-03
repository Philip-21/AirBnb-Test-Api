package handlers

import (
	"airbnb/middleware"
	"airbnb/models"
	"airbnb/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BookingHandlers struct {
	DbRepo *repository.BookingRepo
}

func NewBookingHandlers(repo *repository.BookingRepo) *BookingHandlers {
	return &BookingHandlers{
		DbRepo: repo,
	}
}

// @Tags		   Bookings
// @Summary		   Book Property
// @Description    A User Books a property or apartment
// @Success        200   "successfully booked"
// @Param           propertyid path string true "ID"
// @Router         /user/booking/{propertyid} [post]
// @Param          Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (h *BookingHandlers) CreateBooking(ctx *gin.Context) {
	idParam := ctx.Param("propertyid")
	propertyID, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid property ID"})
		return
	}
	user, err := middleware.GetUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = h.DbRepo.CreateBooking(ctx, &models.Booking{
		UserID:     user.ID,
		PropertyID: propertyID,
		Status:     models.Pending,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "successfully booked"})
}

// @Tags		   Bookings
// @Summary		   Get Bookings
// @Description    A User gets his list of bookings
// @Success        200 {object} []models.UserGetBooking
// @Router         /user/booking [get]
// @Param          Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (h *BookingHandlers) GetUserBookings(ctx *gin.Context) {
	user, err := middleware.GetUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	bookings, err := h.DbRepo.GetUserBookings(ctx, user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, models.GetUserBookings{Bookings: bookings})
}

// @Tags		   Bookings
// @Summary		   Get Bookings
// @Description    A User gets a booking
// @Success        200 {object} models.UserGetBooking
// @Param          bookingid path string true "ID"
// @Router         /user/booking/{bookingid} [get]
// @Param          Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (h *BookingHandlers) GetUserBookingByID(ctx *gin.Context) {
	bookingIDParam := ctx.Param("bookingid")
	bookingID, err := uuid.Parse(bookingIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking ID"})
		return
	}
	_, err = middleware.GetUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	booking, err := h.DbRepo.GetUserBookingByID(ctx, bookingID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, booking)
}

// @Tags		   Bookings
// @Summary		   Confirm Bookings
// @Description    A Property owner or a user can cancel a booking
// @Success        200 "booking confirmed"
// @Param          bookingid path string true "ID"
// @Router         /cancel/booking/{bookingid} [delete]
func (h *BookingHandlers) CancelBooking(ctx *gin.Context) {
	idParam := ctx.Param("bookingid")
	bookingID, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking ID"})
		return
	}

	if err := h.DbRepo.CancelBooking(ctx, bookingID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "booking cancelled"})
}

// @Tags		   Bookings
// @Summary		   Get Bookings
// @Description    A Property owner gets all  booking
// @Success        200 {object} []models.PropertyBooking
// @Router         /owner/booking/all [get]
// @Param          Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (h *BookingHandlers) GetPropertyBookings(ctx *gin.Context) {
	owner, err := middleware.GetPropertyOwner(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	bookings, err := h.DbRepo.GetPropertyBookings(ctx, owner.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, models.GetPropertyBookings{Bookings: bookings})
}

// @Tags		   Bookings
// @Summary		   Get Bookings
// @Description    A Property owner gets a particular  bookings data
// @Success        200 {object} models.PropertyBooking
// @Param          bookingid path string true "ID"
// @Router         /owner/booking/{bookingid} [get]
// @Param          Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (h *BookingHandlers) GetPropertyBookingByID(ctx *gin.Context) {
	bookingIDParam := ctx.Param("bookingid")
	bookingID, err := uuid.Parse(bookingIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid property ID"})
		return
	}

	booking, err := h.DbRepo.GetPropertyBookingByID(ctx, bookingID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, booking)
}

// @Tags		   Bookings
// @Summary		   Confirm Bookings
// @Description    A Property owner confirms a booking
// @Success        200 "booking confirmed"
// @Param          bookingid path string true "ID"
// @Router         /owner/booking/{bookingid} [put]
// @Param          Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (h *BookingHandlers) ConfirmBooking(ctx *gin.Context) {
	idParam := ctx.Param("bookingid")
	bookingID, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking ID"})
		return
	}
	_, err = middleware.GetPropertyOwner(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.DbRepo.ConfirmBooking(ctx, bookingID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "booking confirmed"})
}
