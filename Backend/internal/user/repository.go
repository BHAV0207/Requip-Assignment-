package user

import (
	"context"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IRepository interface {
	Create(ctx context.Context, user *User) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByPAN(ctx context.Context, PAN string) (bool, error)
	ExistsByAadhaar(ctx context.Context, Aadhaar string) (bool, error)
	ExistsByPrimaryMobile(ctx context.Context, mobile string) (bool, error)
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetAll(ctx context.Context, offset, limit int) ([]User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, user *User) error
	Count(ctx context.Context) (int64, error)
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) IRepository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, user *User) error {
	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			msg := mysqlErr.Message
			if strings.Contains(msg, "idx_users_email") {
				return ErrEmailAlreadyExists
			}
			if strings.Contains(msg, "idx_users_primary_mobile") {
				return ErrMobileAlreadyExists
			}
			if strings.Contains(msg, "idx_users_pan") {
				return ErrPANAlreadyExists
			}
			if strings.Contains(msg, "idx_users_aadhaar") {
				return ErrAadhaarAlreadyExists
			}
		}
		return err
	}
	return nil
}

func (r *Repository) ExistsByEmail(ctx context.Context, email string) (bool, error) {

	var count int64

	err := r.db.WithContext(ctx).Unscoped().
		Model(&User{}).
		Where("email = ?", email).
		Count(&count).
		Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *Repository) ExistsByPAN(ctx context.Context, PAN string) (bool, error) {

	var count int64

	err := r.db.WithContext(ctx).Unscoped().
		Model(&User{}).
		Where("PAN = ?", PAN).
		Count(&count).
		Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *Repository) ExistsByAadhaar(ctx context.Context, Aadhaar string) (bool, error) {

	var count int64

	err := r.db.WithContext(ctx).Unscoped().
		Model(&User{}).
		Where("Aadhaar = ?", Aadhaar).
		Count(&count).
		Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *Repository) ExistsByPrimaryMobile(ctx context.Context, mobile string) (bool, error) {

	var count int64

	err := r.db.WithContext(ctx).Unscoped().
		Model(&User{}).
		Where("primary_mobile = ?", mobile).
		Count(&count).
		Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {

	var user User

	err := r.db.WithContext(ctx).
		First(&user, "id = ?", id).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetAll(ctx context.Context, offset, limit int) ([]User, error) {

	var users []User

	err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&users).
		Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Repository) Update(ctx context.Context, user *User) error {

	return r.db.WithContext(ctx).Save(user).Error

}

func (r *Repository) Delete(ctx context.Context, user *User) error {

	return r.db.WithContext(ctx).Delete(user).Error

}

func (r *Repository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&User{}).Count(&count).Error
	return count, err
}
