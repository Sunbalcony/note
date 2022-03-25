package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"note/model"
	"note/util"
	"strings"
)

func ProcessHandleRedis(ctx *gin.Context) {
	fmt.Println(strings.Split(ctx.Request.URL.EscapedPath(), "/")[1])
	msg := &model.Message{}
	if ctx.Request.Method == http.MethodPost {
		fmt.Println(ctx.Request.FormValue("text"))
		msg.Tag = strings.Split(ctx.Request.URL.EscapedPath(), "/")[1]
		msg.Text = ctx.Request.FormValue("text")
		model.Rds.SetKey(msg.Tag, msg.Text)

	}
	if ctx.Request.Method == http.MethodGet {
		if ctx.Request.URL.Path != "/favicon.ico" {
			value, err := model.Rds.GetKey(strings.Split(ctx.Request.URL.EscapedPath(), "/")[1])
			if err != nil {
				ctx.HTML(http.StatusOK, "index.html", gin.H{"title": "免登录 可分享 实时保存的记事本", "text": nil})
				return
			}
			ctx.HTML(http.StatusOK, "index.html", gin.H{"title": "免登录 可分享 实时保存的记事本", "text": value})
			return
		}

	}

}

// CreateRedis Create Api POST请求传递参数Text
func CreateRedis(ctx *gin.Context) {
	msg := &model.Message{}
	ctx.Bind(msg)
	randomTag := util.RandChar(viper.GetInt("note.keylength"))
	if msg.Text != "" {
		exist := model.Rds.KeyExist(randomTag)
		if exist == 1 {
			ctx.String(http.StatusInternalServerError, "can not use CREATE API to create ")
			return
		}
		err := model.Rds.SetKey(randomTag, msg.Text)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "CreateNotebookTextFailed")
			return
		}
		ctx.String(http.StatusOK, fmt.Sprintf("http://%s/%s", ctx.Request.Host, randomTag))
		return

	}
	ctx.String(http.StatusInternalServerError, "text content length must >1")
	return

}

// UpdateRedis Update POST请求传递参数Tag和Text
func UpdateRedis(ctx *gin.Context) {
	msg := &model.Message{}
	ctx.Bind(msg)
	exist := model.Rds.KeyExist(msg.Tag)
	if exist == 0 {
		ctx.String(http.StatusInternalServerError, "will to update message key is not exist")
		return
	}
	err := model.Rds.SetKey(msg.Tag, msg.Text)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "UpdateNotebookTextFailed")
		return
	}
	ctx.String(http.StatusOK, "UpdateNotebookTextSuccess")
	return

}
