package main

import (
	"embed"
	"fmt"
	"html/template"
	"note/model"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

//go:embed static/*
var staticFS embed.FS

func main() {
	r := gin.Default()
	tmpl := template.Must(template.New("").ParseFS(staticFS, "static/*.html"))
	r.SetHTMLTemplate(tmpl)
	r = NewRoutes(r)
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

	switch {
	case class == 0:
		//初始化数据库连接
		db, err := model.NewDBEngine()
		if err != nil {
			panic(errors.Wrap(err, "初始化数据库链接异常"))
		}
		model.DBEngine = db
		return
	case class == 1:
		apiClient := model.NewRedisApi()
		model.Rds = apiClient
		return
	default:
		panic("note.type must be 0 or 1")

	}

}
