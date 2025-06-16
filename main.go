package main

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"html/template"
	"io/fs"
	"net/http"
	"note/model"
)

//go:embed static/*
var staticFS embed.FS

func main() {
	r := gin.Default()
	tmpl := template.Must(template.New("").ParseFS(staticFS, "static/*.html"))
	r.SetHTMLTemplate(tmpl)
	staticContent, _ := fs.Sub(staticFS, "static")
	r = NewRoutes(r)
	r.GET("/static/:file", func(c *gin.Context) {
		file := c.Param("file")
		c.FileFromFS(file, http.FS(staticContent))
	})
	r.GET("/:filename", func(c *gin.Context) {
		filename := c.Param("filename")

		// 检查请求的是否是CSS/JS/ICO文件
		switch filename {
		case "style.css", "script.js", "favicon.ico":
			c.FileFromFS(filename, http.FS(staticContent))
		default:
			c.Next() // 继续路由匹配
		}
	})
	port := viper.GetString("note.serverPort")
	if port != "" {
		_ = r.Run(":" + port)
	}
	panic(r.Run())

}

func init() {
	InitConfig()
	class := viper.GetInt("note.type")
	fmt.Println(class)
	if class == 0 {
		//初始化数据库连接
		db, err := model.NewDBEngine()
		if err != nil {
			panic(errors.Wrap(err, "初始化数据库链接异常"))
		}
		model.DBEngine = db
		return

	}
	if class == 1 {
		apiClient := model.NewRedisApi()
		model.Rds = apiClient
		return

	}
}
