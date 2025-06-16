package main

import (
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

func main() {
	r := gin.Default()
	tmpl := template.Must(template.New("").ParseFS(FS, "static/*.html"))
	r.SetHTMLTemplate(tmpl)
	subStaticFS, _ := fs.Sub(FS, "static")
	r.StaticFS("/static", http.FS(subStaticFS))

	r = NewRoutes(r)
	port := viper.GetString("note.serverPort")
	if port != "" {
		r.Run(":" + port) //监听端口
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
