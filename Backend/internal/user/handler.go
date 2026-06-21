package user

import (
	"net/http"
	"strconv"

	"github.com/BHAV0207/Requip-assignment/internal/response"
	"github.com/gin-gonic/gin"
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
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.Create(req)
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

	user, err := h.service.GetByID(id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "User retrieved successfully", ToUserResponse(user))
}

func (h *Handler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	users, total, err := h.service.GetAll(page, limit)
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
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.Update(id, req)
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

	if err := h.service.Delete(id); err != nil {
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
