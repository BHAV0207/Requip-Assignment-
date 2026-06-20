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
