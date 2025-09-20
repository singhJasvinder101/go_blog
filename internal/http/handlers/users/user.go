package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/singhJasvinder101/go_blog/internal/utils/response"
	"github.com/singhJasvinder101/go_blog/storage/services"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) CreateUserHandler(c *gin.Context) {
	var body struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err))
		return
	}

	id, err := h.service.CreateUser(body.Name, body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
		"user_id": id,
	})
}