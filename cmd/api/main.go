package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/singhJasvinder101/go_blog/internal/config"
	post_handlers "github.com/singhJasvinder101/go_blog/internal/http/handlers/posts"
	user_handlers "github.com/singhJasvinder101/go_blog/internal/http/handlers/users"
	"github.com/singhJasvinder101/go_blog/internal/middleware"
	"github.com/singhJasvinder101/go_blog/internal/utils/jwt"
	"github.com/singhJasvinder101/go_blog/storage/postgres"
	"github.com/singhJasvinder101/go_blog/storage/services"
)

func main() {
	println("Hello")

	// configs setup
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Println("error while loading config:", err.Error())
		panic(err)
	}

	// JWT setup
	jwt.Init(cfg.JwtSecret)
	println("config secret ", cfg.JwtSecret)

	// DB setup
	db, err := postgres.NewPostgres(cfg)
	if err != nil {
		fmt.Println("error while connecting to database:", err.Error())
		panic(err)
	}

	// All tables creation
	db.InitSchema(context.Background())

	// Postgres -> Repo -> Service -> Handler

	// depancies on each other
	// repos initialization
	userRepo := postgres.NewUserRepo(db)
	postRepo := postgres.NewPostRepo(db)
	
	// services initialization
	postService := services.NewPostService(postRepo)
	userService := services.NewUserService(userRepo, postRepo)

	userHandler := user_handlers.NewUserHandler(userService)
	postHandler := post_handlers.NewPostHandler(postService)

	// router setup
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	
	router.POST("/api/users/register", userHandler.RegisterUser)
	router.POST("/api/users/login", userHandler.LoginUser)
	
	// protected routes
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// user routes
		router.POST("/api/users", userHandler.CreateUser)
		router.GET("/api/users/:id", userHandler.GetUserByID)
		router.GET("/api/users/:id/posts", userHandler.GetUserPosts)
		router.POST("/api/users/:user_id/posts/:post_id/comment", userHandler.CreateUserComment)

		// posts routes
		router.POST("/api/posts", postHandler.CreatePost)
		router.GET("/api/posts", postHandler.GetAllPosts)
		router.GET("/api/posts/:id", postHandler.GetPostByID)
	}


	// server setup
	router.Run(fmt.Sprintf(":%d", cfg.HttpServer.Port))
}
