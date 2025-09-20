package post_handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/singhJasvinder101/go_blog/internal/utils/response"
	"github.com/singhJasvinder101/go_blog/storage/services"
)

type PostHandler struct {
	PostService *services.PostService
}

func NewPostHandler(service *services.PostService) *PostHandler {
	return &PostHandler{
		PostService: service,
	}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	ctx := c.Request.Context()
	var body struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description" binding:"required"`
		UserId      int `json:"user_id" binding:"required"`
	}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err))
		return
	}

	id, err := h.PostService.Create(ctx, body.Title, body.Description, body.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "post created successfully",
		"post_id": id,
	})
}

func (h *PostHandler) GetAllPosts(c *gin.Context) {
	ctx := c.Request.Context()

	posts, err := h.PostService.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err))
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
}

func (h *PostHandler) GetPostByID(c *gin.Context){
	ctx := c.Request.Context()
	idParams := c.Param("id")

	postId, err := strconv.Atoi(idParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err))
		return
	}

	post, err := h.PostService.GetByID(ctx, postId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

