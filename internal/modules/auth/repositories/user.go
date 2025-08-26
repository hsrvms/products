package repositories

import (
	"context"
	"fmt"
	"go-starter/internal/modules/auth/models"
	"go-starter/pkg/db"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *db.Database
}

func NewRepository(db *db.Database) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	// Generate UUID if not already set
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	result := r.db.DB.WithContext(ctx).Create(user)
	if result.Error != nil {
		return fmt.Errorf("failed to create user: %w", result.Error)
	}

	return nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	result := r.db.DB.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", result.Error)
	}

	return &user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User

	result := r.db.DB.WithContext(ctx).Where("id = ?", id).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", result.Error)
	}

	return &user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	result := r.db.DB.WithContext(ctx).Save(user)
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}

	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	result := r.db.DB.WithContext(ctx).Delete(&models.User{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	var count int64

	result := r.db.DB.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("failed to check email existence: %w", result.Error)
	}

	return count > 0, nil
}

func (r *UserRepository) ListUsers(ctx context.Context, limit, offset int) ([]*models.User, error) {
	var users []*models.User

	result := r.db.DB.WithContext(ctx).Limit(limit).Offset(offset).Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to list users: %w", result.Error)
	}

	return users, nil
}

func (r *UserRepository) GetUserCount(ctx context.Context) (int64, error) {
	var count int64

	result := r.db.DB.WithContext(ctx).Model(&models.User{}).Count(&count)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to get user count: %w", result.Error)
	}

	return count, nil
}

// Migrate creates the users table using GORM auto-migration
func (r *UserRepository) Migrate(ctx context.Context) error {
	if err := r.db.DB.WithContext(ctx).AutoMigrate(&models.User{}); err != nil {
		return fmt.Errorf("failed to migrate users table: %w", err)
	}

	return nil
}
