package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *Repository) ExistsByEmail(email string) (bool, error) {

	var count int64

	err := r.db.
		Model(&User{}).
		Where("email = ?", email).
		Count(&count).
		Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *Repository) ExistsByPAN(PAN string) (bool, error) {

	var count int64

	err := r.db.
		Model(&User{}).
		Where("PAN = ?", PAN).
		Count(&count).
		Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *Repository) ExistsByAadhaar(Aadhaar string) (bool, error) {

	var count int64	

	err := r.db.
		Model(&User{}).
		Where("Aadhaar = ?", Aadhaar).
		Count(&count).
		Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}



func (r *Repository) FindByID(id uuid.UUID) (*User, error) {

	var user User

	err := r.db.
		First(&user, "id = ?", id).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetAll(offset, limit int) ([]User, error) {

	var users []User

	err := r.db.
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

func (r *Repository) Update(user *User) error {

	return r.db.Save(user).Error

}

func (r *Repository) Delete(user *User) error {

	return r.db.Delete(user).Error

}
