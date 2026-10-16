package user

import (
	"net/http"
	"strconv"

	"github.com/BHAV0207/Requip-assignment/internal/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Handler struct {
	service IService
}

func NewHandler(service IService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, formatValidationError(err))
		return
	}

	user, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "User created successfully", ToUserResponse(user))
}

func (h *Handler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "User retrieved successfully", ToUserResponse(user))
}

func (h *Handler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	users, total, err := h.service.GetAll(c.Request.Context(), page, limit)
	if err != nil {
		h.handleError(c, err)
		return
	}

	userResponses := make([]UserResponse, len(users))
	for i, u := range users {
		userResponses[i] = ToUserResponse(&u)
	}

	response.Success(c, http.StatusOK, "Users retrieved successfully", gin.H{
		"users": userResponses,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *Handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, formatValidationError(err))
		return
	}

	user, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "User updated successfully", ToUserResponse(user))
}

func (h *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		h.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "User deleted successfully", nil)
}

func (h *Handler) handleError(c *gin.Context, err error) {
	switch err {
	case ErrEmailAlreadyExists, ErrPANAlreadyExists, ErrAadhaarAlreadyExists, ErrMobileAlreadyExists:
		response.Error(c, http.StatusConflict, err.Error())
	case ErrUserNotFound:
		response.Error(c, http.StatusNotFound, err.Error())
	default:
		response.Error(c, http.StatusInternalServerError, "internal server error")
	}
}

func formatValidationError(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := e.Field()
			switch e.Tag() {
			case "required":
				return field + " is required"
			case "email":
				return "Invalid email format"
			case "len":
				return field + " must be exactly " + e.Param() + " characters"
			case "min":
				return field + " must be at least " + e.Param() + " characters"
			case "max":
				return field + " must be at most " + e.Param() + " characters"
			}
		}
	}
	return "Invalid request body"
}
