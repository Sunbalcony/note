package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"net/http"
	"note/model"
	"note/util"
	"strings"
)

type NoteApi struct {
	model.Message
	db *gorm.DB
}

func NewNoteApi() *NoteApi {
	return &NoteApi{db: model.DBEngine}

}

func ProcessHandleMysql(ctx *gin.Context) {
	fmt.Println(strings.Split(ctx.Request.URL.EscapedPath(), "/")[1])
	msg := &model.Message{}
	if ctx.Request.Method == http.MethodPost {
		fmt.Println(ctx.Request.FormValue("text"))
		//dataSet[ctx.Request.RequestURI] = ctx.PostForm("text")
		if model.DBEngine.Where("tag = ?", strings.Split(ctx.Request.URL.EscapedPath(), "/")[1]).First(&msg).RecordNotFound() == true {
			msg.Tag = strings.Split(ctx.Request.URL.EscapedPath(), "/")[1]
			msg.Text = ctx.Request.FormValue("text")
			model.DBEngine.Create(&msg)
		}
		if model.DBEngine.Where("tag = ?", strings.Split(ctx.Request.URL.EscapedPath(), "/")[1]).First(&msg).RecordNotFound() == false {
			model.DBEngine.Model(&msg).Where("tag = ?", strings.Split(ctx.Request.URL.EscapedPath(), "/")[1]).Update("text", ctx.Request.FormValue("text"))
		}

	}
	if ctx.Request.Method == http.MethodGet {
		if ctx.Request.URL.Path != "/favicon.ico" {
			model.DBEngine.Where("tag = ?", strings.Split(ctx.Request.URL.EscapedPath(), "/")[1]).First(&msg)
			ctx.HTML(http.StatusOK, "index.html", gin.H{"title": "免登录 可分享 实时保存的记事本", "text": msg.Text})
		}

	}

}

// Create Api POST请求传递参数Text
func (r *NoteApi) Create(ctx *gin.Context) {
	msg := &model.Message{}
	ctx.Bind(&r.Message)
	randomTag := util.RandChar(viper.GetInt("note.keylength"))
	if r.db.Where("tag = ?", randomTag).First(msg).RecordNotFound() == true {
		r.Message.Tag = randomTag
		if affected := r.db.Create(r.Message).RowsAffected; affected == 1 {
			//如果你配置了nginx代理并启用了HTTPS 下面就填https如果没有就是http
			ctx.String(http.StatusOK, fmt.Sprintf("http://%s/%s", ctx.Request.Host, randomTag))
			return

		}
		ctx.String(http.StatusInternalServerError, "CreateNotebookTextFailed")
		return

	}

}

// Update POST请求传递参数Tag和Text
func (r *NoteApi) Update(ctx *gin.Context) {
	msg := &model.Message{}
	ctx.Bind(&r.Message)
	if r.db.Where("tag = ?", r.Message.Tag).First(msg).RecordNotFound() == false {
		if affected := r.db.Model(r.Message).Where("tag = ?", r.Message.Tag).Update("text", &r.Message.Text).RowsAffected; affected != 0 {
			ctx.String(http.StatusOK, "UpdateNotebookTextSuccess")
			return
		}
		ctx.String(http.StatusOK, "UpdateNotebookTextNotChanged")
		return
	}
	ctx.String(http.StatusInternalServerError, "UpdateNotebookTextFailed")

}
