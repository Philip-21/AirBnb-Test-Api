package repository

import (
	"airbnb/models"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, user *models.User) error {
	if err := r.DB.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.DB.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetAllUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	if err := r.DB.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepo) UpdateUser(ctx context.Context, user *models.User) error {
	if err := r.DB.WithContext(ctx).Save(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := r.DB.WithContext(ctx).Where("id = ?", id).Delete(&models.User{}).Error; err != nil {
		return err
	}
	return nil
}
