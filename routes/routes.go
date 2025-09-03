package routes

import (
	_ "airbnb/docs"
	"airbnb/handlers"
	"airbnb/middleware"
	"airbnb/repository"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Routes(
	userRepo *repository.UserRepo,
	propertyRepo *repository.PropertyRepo,
	propertyHandlers *handlers.PropertyHandlers,
	userHandlers *handlers.UserHandlers,
	bookingHandlers *handlers.BookingHandlers,
) *gin.Engine {
	router := gin.Default()

	router.Use(gin.Logger())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/property/all", propertyHandlers.GetProperties)
	router.DELETE("/cancel/booking/:bookingid", bookingHandlers.CancelBooking)

	router.POST("/property/owner/signup", propertyHandlers.CreatePropertyOwner)
	router.POST("/property/owner/login", propertyHandlers.LoginPropertyOwner)
	userRoutes := router.Group("/user")
	{
		userRoutes.POST("/signup", userHandlers.CreateUser)
		userRoutes.POST("/login", userHandlers.LoginUser)
	}

	userBookingRoutes := router.Group("/user")
	userBookingRoutes.Use(middleware.AuthUser(userRepo))
	{
		userBookingRoutes.POST("/booking/:propertyid", bookingHandlers.CreateBooking)
		userBookingRoutes.GET("/booking", bookingHandlers.GetUserBookings)
		userBookingRoutes.GET("/booking/:bookingid", bookingHandlers.GetUserBookingByID)
	}

	propertyRoutes := router.Group("/property")
	propertyRoutes.Use(middleware.AuthPropertyOwner(propertyRepo))
	{
		propertyRoutes.POST("/create", propertyHandlers.CreateProperty)
		propertyRoutes.GET("/:propertyid", propertyHandlers.GetPropertyByID)
		propertyRoutes.GET("/owner", propertyHandlers.GetAllProperties)
	}
	ownerBookingRoutes := router.Group("/owner/booking")
	ownerBookingRoutes.Use(middleware.AuthPropertyOwner(propertyRepo))
	{
		ownerBookingRoutes.GET("/all", bookingHandlers.GetPropertyBookings)
		ownerBookingRoutes.GET("/:bookingid", bookingHandlers.GetPropertyBookingByID)
		ownerBookingRoutes.PUT("/:bookingid", bookingHandlers.ConfirmBooking)
	}

	return router
}
