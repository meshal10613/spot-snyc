package repository

import (
	"errors"
	"spot-sync/internal/models"

	"gorm.io/gorm"
)

// AuthRepository defines all database operations needed by the auth service.
type AuthRepository interface {
	CreateUser(user *models.User) error
	FindUserByEmail(email string) (*models.User, error)
	FindUserByID(id uint) (*models.User, error)
}

type authRepository struct {
	db *gorm.DB
}

// NewAuthRepository creates a new auth repository instance.
func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *authRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // not found is not an error — caller decides
	}
	return &user, err
}

func (r *authRepository) FindUserByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}
