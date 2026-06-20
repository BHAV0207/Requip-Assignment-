package user

import (
	"net/http"

	"github.com/BHAV0207/Requip-assignment/internal/response"
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

	var req CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.Create(req)

	if err != nil {

		switch err {

		case ErrEmailAlreadyExists,
			ErrPANAlreadyExists,
			ErrAadhaarAlreadyExists:

			response.Error(c, http.StatusConflict, err.Error())

		default:

			response.Error(c, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	response.Success(
		c,
		http.StatusCreated,
		"User created successfully",
		ToUserResponse(user),
	)
}
