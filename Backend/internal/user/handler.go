package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Create(c *gin.Context) {

	var request CreateUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})

		return
	}

	user := User{
		Name: request.Name,

		Email: request.Email,

		PrimaryMobile: request.PrimaryMobile,

		SecondaryMobile: request.SecondaryMobile,

		Aadhaar: request.Aadhaar,

		PAN: request.PAN,

		PlaceOfBirth: request.PlaceOfBirth,

		CurrentAddress: request.CurrentAddress,

		PermanentAddress: request.PermanentAddress,
	}

	err := h.service.Create(&user)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User created successfully",
		"data":    user,
	})
}