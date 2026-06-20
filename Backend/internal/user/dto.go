package user

type CreateUserRequest struct {
	Name             string `json:"name" validate:"required,min=2,max=100"`
	Email            string `json:"email" validate:"required,email"`
	PrimaryMobile    string `json:"primaryMobile" validate:"required,len=10"`
	SecondaryMobile  string `json:"secondaryMobile"`
	Aadhaar          string `json:"aadhaar" validate:"required,len=12"`
	PAN              string `json:"pan" validate:"required,len=10"`
	DateOfBirth      string `json:"dateOfBirth"`
	PlaceOfBirth     string `json:"placeOfBirth"`
	CurrentAddress   string `json:"currentAddress"`
	PermanentAddress string `json:"permanentAddress"`
}

type UserResponse struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Email            string `json:"email"`
	PrimaryMobile    string `json:"primaryMobile"`
	SecondaryMobile  string `json:"secondaryMobile,omitempty"`
	DateOfBirth      string `json:"dateOfBirth"`
	PlaceOfBirth     string `json:"placeOfBirth"`
	CurrentAddress   string `json:"currentAddress"`
	PermanentAddress string `json:"permanentAddress"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
}
