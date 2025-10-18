package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/singhJasvinder101/go_blog/internal/config"
	"github.com/singhJasvinder101/go_blog/internal/middleware"
	"github.com/singhJasvinder101/go_blog/internal/utils/jwt"
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

	app, err := InitializeApplication(cfg)
	if err != nil {
		fmt.Println("error while initializing application:", err.Error())
		panic(err)
	}

	defer app.DB.Close()
	defer app.Redis.Client.Close()

	//redis setup
	ctx := context.Background()
	pong, err := app.Redis.Client.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Could not connect to Redis:", err)
		return
	}
	fmt.Println("Connected to Redis:", pong)

	// all tables creation
	err = app.DB.InitSchema(ctx)
	if err != nil {
		return
	}

	// router setup
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/api/users/register", app.UserHandler.RegisterUser)
	router.POST("/api/users/login", app.UserHandler.LoginUser)

	// protected routes
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// user routes
		router.POST("/api/users", app.UserHandler.CreateUser)
		router.GET("/api/users/:id", app.UserHandler.GetUserByID)
		router.GET("/api/users/:id/posts", app.UserHandler.GetUserPosts)
		router.POST("/api/users/:user_id/posts/:post_id/comment", app.UserHandler.CreateUserComment)

		// posts routes
		router.POST("/api/posts", app.PostHandler.CreatePost)
		router.GET("/api/posts", app.PostHandler.GetAllPosts)
		router.GET("/api/posts/:id", app.PostHandler.GetPostByID)
	}

	// server setup
	router.Run(fmt.Sprintf(":%d", cfg.HttpServer.Port))
}
