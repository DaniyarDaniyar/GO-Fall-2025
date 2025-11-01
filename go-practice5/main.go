package main

import (
	"github.com/DaniyarDaniyar/go-practice5/config"
	"github.com/DaniyarDaniyar/go-practice5/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()
	defer config.DB.Close()

	r := gin.Default()

	r.GET("/users", handlers.GetUsers)

	r.Run(":8080")
}
