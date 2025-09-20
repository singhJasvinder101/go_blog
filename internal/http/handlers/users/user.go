package user_handlers

import (
	"errors"
	"net/http"
	"strconv"

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

func (h *UserHandler) CreateUser(c *gin.Context) {
	ctx := c.Request.Context()

	var body struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err))
		return
	}

	id, err := h.service.Create(ctx, body.Name, body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
		"user_id": id,
	})
}

func (h *UserHandler) GetUserByID(c *gin.Context){
	ctx := c.Request.Context()
	idParams := c.Param("id")

	if idParams == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(errors.New("id parameter is required")))
		return
	}

	id, err := strconv.Atoi(idParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(errors.New("invalid id parameter")))
		return
	}

	user, err := h.service.GetByID(ctx, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (h *UserHandler) GetUserPosts(c *gin.Context){
	ctx := c.Request.Context()
	user_id := c.Param("id")

	if user_id == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(errors.New("id parameter is required")))
		return
	}
	
	uid, err := strconv.Atoi(user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(errors.New("invalid id parameter")))
		return
	}

	posts, err := h.service.GetByUID(ctx, uid)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err))
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})

}

