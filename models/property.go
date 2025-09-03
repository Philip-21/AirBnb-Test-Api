package models

import "github.com/google/uuid"

type CreatePropertyOwner struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GetPropertyOwner struct {
	OwnerID uuid.UUID `json:"owner_id"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
}
type LoginPropertyOwner struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type CreateProperty struct {
	PropertyName string `json:"property_name"`
	Description  string `json:"description"`
	Price        int64  `json:"price"`
}

type GetProperty struct {
	PropertyID    uuid.UUID        `json:"property_id"`
	PropertyName  string           `json:"property_name"`
	Description   string           `json:"description"`
	Price         int64            `json:"price"`
	PropertyOwner GetPropertyOwner `json:"property_owner"`
}

type GetAllProperties struct {
	Properties []GetProperty `json:"properties"`
}
