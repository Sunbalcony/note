package main

import "github.com/gin-gonic/gin"

// NewRoutes 路由
func NewRoutes(r *gin.Engine) *gin.Engine {
	api := NewNoteApi()
	r.GET("/", index)
	r.POST("/:id", processHandle)
	r.GET("/:id", processHandle)
	apiRouter := r.Group("/api")
	apiRouter.POST("/create", api.create)
	apiRouter.POST("/update", api.update)

	r.Static("./static", "./static")
	r.LoadHTMLGlob("./home/*")
	return r

}
