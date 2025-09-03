package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID      `gorm:"column:id;type:uuid;primaryKey"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

type User struct {
	BaseModel
	Name     string `gorm:"size:100;not null"`
	Email    string `gorm:"uniqueIndex;size:100;not null"`
	Password string `gorm:"size:255;not null"`
	Role     string `gorm:"size:100;not null"`
}

type PropertyOwner struct {
	BaseModel
	Name       string     `gorm:"size:100;not null"`
	Email      string     `gorm:"uniqueIndex;size:100;not null"`
	Password   string     `gorm:"size:255;not null"`
	Role       string     `gorm:"size:100;not null"`
	Properties []Property `gorm:"foreignKey:OwnerID"` // one-to-many relation
}

type Property struct {
	BaseModel
	Name        string        `gorm:"size:100;not null"`
	Description string        `gorm:"size:500"`
	Price       int64         `gorm:"not null"`
	Location    string        `gorm:"not null"`
	OwnerID     uuid.UUID     `gorm:"type:uuid;not null"`               // foreign key
	Owner       PropertyOwner `gorm:"foreignKey:OwnerID;references:ID"` // GORM association
}

type Booking struct {
	BaseModel
	UserID     uuid.UUID `gorm:"type:uuid;not null"`
	PropertyID uuid.UUID `gorm:"type:uuid;not null"`
	CheckIn    string    `gorm:"default:null"`
	CheckOut   string    `gorm:"default:null"`
	Status     string    `gorm:"size:50;not null"` // pending, confirmed
	User       User      `gorm:"foreignKey:UserID;references:ID"`
	Property   Property  `gorm:"foreignKey:PropertyID;references:ID"`
}

const (
	UserRole     = "user"
	PropertyRole = "property_owner"
	Pending      = "pending"
	Confirmed    = "confirmed"
	Cancelled    = "cancelled"
)

func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
    if b.ID == uuid.Nil {
        b.ID = uuid.New()
    }
    return
}
