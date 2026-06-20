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

func (r *Repository) FindByEmail(email string) (*User, error) {

	var user User

	err := r.db.
		Where("email = ?", email).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) FindByPAN(PAN string) (*User, error) {

	var user User

	err := r.db.
		Where("PAN = ?", PAN).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) FindByAadhaar(Aadhaar string) (*User, error) {

	var user User

	err := r.db.
		Where("Aadhaar = ?", Aadhaar).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
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
