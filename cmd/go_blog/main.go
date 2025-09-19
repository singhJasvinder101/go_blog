package main

import (
	"fmt"

	"github.com/singhJasvinder101/go_blog/internal/config"
)


func main()  {
	println("Hello")

	_, err := config.NewConfig()

	if err != nil {
		fmt.Println("error while loading config:", err.Error())
		panic(err)
	}

	// configs setup
	// DB setup
	// router setup
	// server setup
	
}