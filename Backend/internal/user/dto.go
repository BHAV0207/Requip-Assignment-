package user

type CreateUserRequest struct {
	Name             string `json:"name" binding:"required,min=2,max=100"`
	Email            string `json:"email" binding:"required,email"`
	PrimaryMobile    string `json:"primaryMobile" binding:"required,len=10"`
	SecondaryMobile  string `json:"secondaryMobile"`
	Aadhaar          string `json:"aadhaar" binding:"required,len=12"`
	PAN              string `json:"pan" binding:"required,len=10"`
	DateOfBirth      string `json:"dateOfBirth" binding:"required"`
	PlaceOfBirth     string `json:"placeOfBirth" binding:"required"`
	CurrentAddress   string `json:"currentAddress" binding:"required"`
	PermanentAddress string `json:"permanentAddress" binding:"required"`
}

type UpdateUserRequest struct {
	Name             *string `json:"name" validate:"omitempty,min=2,max=100"`
	Email            *string `json:"email" validate:"omitempty,email"`
	PrimaryMobile    *string `json:"primaryMobile" validate:"omitempty,len=10"`
	SecondaryMobile  *string `json:"secondaryMobile"`
	Aadhaar          *string `json:"aadhaar" validate:"omitempty,len=12"`
	PAN              *string `json:"pan" validate:"omitempty,len=10"`
	DateOfBirth      *string `json:"dateOfBirth"`
	PlaceOfBirth     *string `json:"placeOfBirth"`
	CurrentAddress   *string `json:"currentAddress"`
	PermanentAddress *string `json:"permanentAddress"`
}

type UserResponse struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Email            string `json:"email"`
	PrimaryMobile    string `json:"primaryMobile"`
	SecondaryMobile  string `json:"secondaryMobile,omitempty"`
	Aadhaar          string `json:"aadhaar"`
	PAN              string `json:"pan"`
	DateOfBirth      string `json:"dateOfBirth"`
	PlaceOfBirth     string `json:"placeOfBirth"`
	CurrentAddress   string `json:"currentAddress"`
	PermanentAddress string `json:"permanentAddress"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
}
