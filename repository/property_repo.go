package repository

import (
	"airbnb/models"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PropertyRepo struct {
	DB *gorm.DB
}

func NewPropertyRepo(db *gorm.DB) *PropertyRepo {
	return &PropertyRepo{
		DB: db,
	}
}

func (r *PropertyRepo) CreateProperty(ctx context.Context, property *models.Property) error {
	if err := r.DB.WithContext(ctx).Create(property).Error; err != nil {
		return fmt.Errorf("failed to create property: %w", err)
	}
	return nil
}

func (r *PropertyRepo) GetPropertyByID(ctx context.Context, id uuid.UUID) (*models.Property, error) {
	var property models.Property
	if err := r.DB.WithContext(ctx).Preload("Owner").First(&property, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch property: %w", err)
	}
	return &property, nil
}

func (r *PropertyRepo) GetAllProperties(ctx context.Context, ownerID uuid.UUID) ([]models.Property, error) {
	var properties []models.Property
	if err := r.DB.WithContext(ctx).Where("owner_id = ?", ownerID).Preload("Owner").Find(&properties).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch properties: %w", err)
	}
	return properties, nil
}

func (r *PropertyRepo) GetProperties(ctx context.Context) ([]models.Property, error) {
	var properties []models.Property
	if err := r.DB.WithContext(ctx).Preload("Owner").Find(&properties).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch properties: %w", err)
	}
	return properties, nil
}

func (r *PropertyRepo) UpdateProperty(ctx context.Context, property *models.Property) error {
	property.UpdatedAt = time.Now()
	if err := r.DB.WithContext(ctx).Save(property).Error; err != nil {
		return fmt.Errorf("failed to update property: %w", err)
	}
	return nil
}

func (r *PropertyRepo) DeleteProperty(ctx context.Context, id uuid.UUID) error {
	if err := r.DB.WithContext(ctx).Delete(&models.Property{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete property: %w", err)
	}
	return nil
}

func (r *PropertyRepo) GetPropertiesByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]models.Property, error) {
	var properties []models.Property
	if err := r.DB.WithContext(ctx).Where("owner_id = ?", ownerID).Preload("Owner").Find(&properties).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch properties for owner: %w", err)
	}
	return properties, nil
}

func (r *PropertyRepo) CreatePropertyOwner(ctx context.Context, owner *models.PropertyOwner) error {
	if err := r.DB.WithContext(ctx).Create(owner).Error; err != nil {
		return fmt.Errorf("failed to create property owner: %w", err)
	}
	return nil
}

func (r *PropertyRepo) GetPropertyOwnerByID(ctx context.Context, ownerID uuid.UUID) (*models.PropertyOwner, error) {
	var owner models.PropertyOwner
	if err := r.DB.WithContext(ctx).First(&owner, "id = ?", ownerID).Error; err != nil {
		return nil, fmt.Errorf("property owner not found: %w", err)
	}
	return &owner, nil
}

func (r *PropertyRepo) GetPropertyOwnerByEmail(ctx context.Context, email string) (*models.PropertyOwner, error) {
	var owner models.PropertyOwner
	if err := r.DB.WithContext(ctx).Where("email = ?", email).First(&owner).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("property owner not found")
		}
		return nil, err
	}
	return &owner, nil
}
