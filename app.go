package main

import (
	"github.com/gin-gonic/gin"
)

type App struct {
	Gin *gin.Engine
}

type AppOptions struct {
}

func NewApp(opts ...AppOptions) *App {
	//options := AppOptions{}
	//for _, i := range opts {
	//options = i
	//break
	//}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	app := &App{
		Gin: r,
	}

	return app
}

func (app *App) Run() {
}
