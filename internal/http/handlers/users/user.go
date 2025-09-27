package user_handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/singhJasvinder101/go_blog/internal/utils/hash"
	"github.com/singhJasvinder101/go_blog/internal/utils/jwt"
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

func (h *UserHandler) RegisterUser(c *gin.Context) {
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

	hashedPassword, err := hash.HashPassword(body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err))
		return
	}

    createdUser, err := h.service.Create(ctx, body.Name, body.Email, hashedPassword)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

	token, err := jwt.GenerateToken(createdUser.ID, createdUser.Name, createdUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err))
		return
	}

	c.SetCookie("auth_token", token, 3600*24, "/", "localhost", false, true)
	
	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
		"user":    createdUser,
		"token":   token,
	})
	
}

func (h *UserHandler) LoginUser(c *gin.Context){
	ctx := c.Request.Context()
	
	var body struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err))
		return
	}

	user, err := h.service.GetByEmail(ctx, body.Email)
	if err != nil {
		// c.JSON(http.StatusUnauthorized, response.ErrorResponse(fmt.Errorf("invalid email or password %w", err)))
		c.JSON(http.StatusUnauthorized, response.ErrorResponse(errors.New("invalid email or password")))
		return
	}

	if body.Password == "" || user.PasswordHash == "" {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse(errors.New("invalid email or password")))
		return
	}

	if err := hash.ComparePassword(user.PasswordHash, body.Password); err != nil {
		println("hello2")
    	c.JSON(http.StatusUnauthorized, response.ErrorResponse(fmt.Errorf("invalid email or password %w", err)))
		// c.JSON(http.StatusUnauthorized, response.ErrorResponse(errors.New("invalid email or password")))
		return
	}
	token, err := jwt.GenerateToken(user.ID, user.Name, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err))
		return
	}
	
	c.SetCookie("auth_token", token, 3600*24, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"user":    user,
		"token":   token,
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


func (h * UserHandler) CreateUserComment(c *gin.Context){
	ctx := c.Request.Context()
	id_param1 := c.Param("user_id")
	id_param2 := c.Param("post_id")

	user_id, err1 := strconv.Atoi(id_param1)
	if err1 != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(errors.New("invalid user id parameter")))
		return
	}

	post_id, err2 := strconv.Atoi(id_param2)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(errors.New("invalid post id parameter")))
		return
	}

	var body struct {
		Content string `json:"content" binding:"required"`
	}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err))
		return
	}

	comment, err := h.service.CreateComment(ctx, body.Content, user_id, post_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"comment": comment,
	})
}

