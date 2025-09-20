package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/singhJasvinder101/go_blog/internal/config"
	post_handlers "github.com/singhJasvinder101/go_blog/internal/http/handlers/posts"
	user_handlers "github.com/singhJasvinder101/go_blog/internal/http/handlers/users"
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

	// DB setup
	db, err := postgres.NewPostgres(cfg)
	if err != nil {
		fmt.Println("error while connecting to database:", err.Error())
		panic(err)
	}

	// All tables creation
	db.InitSchema(context.Background())

	// Postgres -> Repo -> Service -> Handler

	// for users
	userRepo := postgres.NewUserRepo(db)
	userService := services.NewUserService(userRepo)
	userHandler := user_handlers.NewUserHandler(userService)

	// for posts
	postRepo := postgres.NewPostRepo(db)
	postService := services.NewPostService(postRepo)
	postHandler := post_handlers.NewPostHandler(postService)

	// router setup
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/api/users", userHandler.CreateUser)
	router.GET("/api/users/:id", userHandler.GetUserByID)
	router.GET("/api/users/:id/posts", userHandler.GetUserPosts)

	router.POST("/api/posts", postHandler.CreatePost)

	// server setup
	router.Run(fmt.Sprintf(":%d", cfg.HttpServer.Port))
}
