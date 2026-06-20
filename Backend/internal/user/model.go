package user

import (
	"time"

	"github.com/BHAV0207/Requip-assignment/internal/common"
)

type User struct {
	common.BaseModel

	Name string `gorm:"size:100;not null"`

	Email string `gorm:"size:255;uniqueIndex;not null"`

	PrimaryMobile string `gorm:"size:15;uniqueIndex;not null"`

	SecondaryMobile string `gorm:"size:15"`

	Aadhaar string `gorm:"size:12;uniqueIndex;not null"`

	PAN string `gorm:"size:10;uniqueIndex;not null"`

	DateOfBirth *time.Time

	PlaceOfBirth string

	CurrentAddress string

	PermanentAddress string
}

func NewUser(req CreateUserRequest) *User {
	return &User{
		Name:             req.Name,
		Email:            req.Email,
		PrimaryMobile:    req.PrimaryMobile,
		SecondaryMobile:  req.SecondaryMobile,
		Aadhaar:          req.Aadhaar,
		PAN:              req.PAN,
		PlaceOfBirth:     req.PlaceOfBirth,
		CurrentAddress:   req.CurrentAddress,
		PermanentAddress: req.PermanentAddress,
	}
}

func ToUserResponse(user *User) UserResponse {
	return UserResponse{
		ID:                user.ID.String(),
		Name:              user.Name,
		Email:             user.Email,
		PrimaryMobile:     user.PrimaryMobile,
		SecondaryMobile:   user.SecondaryMobile,
		DateOfBirth:       user.DateOfBirth.Format("2006-01-02"),
		PlaceOfBirth:      user.PlaceOfBirth,
		CurrentAddress:    user.CurrentAddress,
		PermanentAddress:  user.PermanentAddress,
		CreatedAt:         user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         user.UpdatedAt.Format(time.RFC3339),
	}
}
