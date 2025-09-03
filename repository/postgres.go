package repository

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectToDB(dsn string) (*gorm.DB, error) {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Println("Error in connection", err)
		return nil, err
	}
	// err = db.AutoMigrate(&models.User{}, &models.PropertyOwner{}, &models.Property{}, &models.Booking{})
	// if err != nil {
	// 	log.Printf("unable to migrate data models : %s", err)
	// 	return nil, err
	// }

	return db, nil
}
