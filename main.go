package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"math/rand"
	"net/http"
	"note/config"
	"strings"
	"time"
)

var DBEngine *gorm.DB

type Message struct {
	Tag  string `json:"tag"`
	Text string `json:"text"`
}
type NoteApi struct {
	Message
	db *gorm.DB
}

const char = "abcdefghijklmnopqrstuvwxyz0123456789"

func RandChar(size int) string {
	rand.NewSource(time.Now().UnixNano())
	var s bytes.Buffer
	for i := 0; i < size; i++ {
		s.WriteByte(char[rand.Int63()%int64(len(char))])
	}
	return s.String()
}
func index(ctx *gin.Context) {
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Header("Pragma", "no-cache")
	ctx.Header("Expires", "0")
	ctx.Redirect(http.StatusMovedPermanently, ctx.Request.URL.EscapedPath()+RandChar(viper.GetInt("note.keylength"))) //url path长度

}

func processHandle(ctx *gin.Context) {
	fmt.Println(strings.Split(ctx.Request.URL.EscapedPath(), "/")[1])
	msg := &Message{}
	if ctx.Request.Method == http.MethodPost {
		fmt.Println(ctx.Request.FormValue("text"))
		//dataSet[ctx.Request.RequestURI] = ctx.PostForm("text")
		if DBEngine.Where("tag = ?", strings.Split(ctx.Request.URL.EscapedPath(), "/")[1]).First(&msg).RecordNotFound() == true {
			msg.Tag = strings.Split(ctx.Request.URL.EscapedPath(), "/")[1]
			msg.Text = ctx.Request.FormValue("text")
			DBEngine.Create(&msg)
		}
		if DBEngine.Where("tag = ?", strings.Split(ctx.Request.URL.EscapedPath(), "/")[1]).First(&msg).RecordNotFound() == false {
			DBEngine.Model(&msg).Where("tag = ?", strings.Split(ctx.Request.URL.EscapedPath(), "/")[1]).Update("text", ctx.Request.FormValue("text"))
		}

	}
	if ctx.Request.Method == http.MethodGet {
		if ctx.Request.URL.Path != "/favicon.ico" {
			DBEngine.Where("tag = ?", strings.Split(ctx.Request.URL.EscapedPath(), "/")[1]).First(&msg)
			ctx.HTML(http.StatusOK, "index.html", gin.H{"title": "免登录 可分享 实时保存的记事本", "text": msg.Text})
		}

	}

}
func NewNoteApi() *NoteApi {
	return &NoteApi{db: DBEngine}

}

// Api POST请求传递参数Text
func (r *NoteApi) create(ctx *gin.Context) {
	msg := &Message{}
	ctx.Bind(&r.Message)
	randomTag := RandChar(10)
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

//POST请求传递参数Tag和Text
func (r *NoteApi) update(ctx *gin.Context) {
	msg := &Message{}
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
func init() {

	//初始化配置文件信息
	config.InitConfig()

	//初始化数据库连接
	db, err := NewDBEngine()
	if err != nil {

		panic(errors.Wrap(err,"初始化数据库链接异常"))
	}
	DBEngine = db


}

func main() {
	r := gin.Default()
	r = NewRoutes(r)
	port := viper.GetString("note.serverPort")
	if port != "" {
		r.Run(":" + port) //监听端口
	}
	panic(r.Run())

}
