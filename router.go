package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"note/service"
	"note/util"
)

func index(ctx *gin.Context) {
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Header("Pragma", "no-cache")
	ctx.Header("Expires", "0")
	ctx.Redirect(http.StatusMovedPermanently, ctx.Request.URL.EscapedPath()+util.RandChar(viper.GetInt("note.keylength"))) //url path长度

}

//go:embed static
var FS embed.FS

// NewRoutes 路由
func NewRoutes(r *gin.Engine) *gin.Engine {

	class := viper.GetInt("note.type")
	if class == 0 {
		api := service.NewNoteApi()
		r.GET("/", index)
		r.POST("/:id", service.ProcessHandleMysql)
		r.GET("/:id", service.ProcessHandleMysql)
		r.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", nil)
		})

		apiRouter := r.Group("/api")
		apiRouter.POST("/create", api.Create)
		apiRouter.POST("/update", api.Update)

		return r
	}
	if class == 1 {
		r.GET("/", index)
		r.POST("/:id", service.ProcessHandleRedis)
		r.GET("/:id", service.ProcessHandleRedis)
		r.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", nil)
		})

		apiRouter := r.Group("/api")
		apiRouter.POST("/create", service.CreateRedis)
		apiRouter.POST("/update", service.UpdateRedis)
		return r
	}
	return r

}
