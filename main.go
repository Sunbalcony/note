package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"math/rand"
	"net/http"
	"net/url"
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

const char = "abcdefghijklmnopqrstuvwxyz0123456789!#&*"

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
	ctx.Redirect(http.StatusMovedPermanently, ctx.Request.URL.EscapedPath()+RandChar(10)) //url path长度

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

// Api POST请求传递参数Tag和Text
func (r *NoteApi) create(ctx *gin.Context) {
	msg := &Message{}
	ctx.Bind(&r.Message)
	randomTag := RandChar(10)
	if r.db.Where("tag = ?", randomTag).First(msg).RecordNotFound() == true {
		r.Message.Tag = randomTag
		if affected := r.db.Create(r.Message).RowsAffected; affected == 1 {
			ctx.String(http.StatusOK, fmt.Sprintf("http://%s/%s", ctx.Request.Host, randomTag))
			return

		}
		ctx.String(http.StatusInternalServerError, "CreateNotebookTextFailed")
		return

	}

}

func (r *NoteApi) update(ctx *gin.Context) {
	msg := &Message{}
	ctx.Bind(&r.Message)
	fmt.Println(r.Message)
	if r.db.Where("tag = ?", r.Message.Tag).First(msg).RecordNotFound() == false {
		if affected := r.db.Model(r.Message).Where("tag = ?", r.Message.Tag).Update("text", &r.Message.Text).RowsAffected; affected != 0 {
			ctx.String(http.StatusOK, "UpdateNotebookTextSuccess")
			return
		}
		ctx.String(http.StatusOK, "UpdateNotebookTextNotChanged")
		return
	}
	ctx.String(http.StatusInternalServerError,"UpdateNotebookTextFailed")

}
func init() {
	db, err := NewDBEngine()
	if err != nil {
		panic(err.Error())
	}
	DBEngine = db

}
func NewDBEngine() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=%s",
		"user",               //用户名
		"password",          //密码
		"mysql.low.im:10073", //db地址
		"notes",              //库名，要先建库
		"utf8",
		true,
		url.QueryEscape("Asia/Shanghai"),
	))
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(30)
	//初始化表结构
	db.AutoMigrate(&Message{})
	return db, nil

}

func main() {
	api := NewNoteApi()
	r := gin.Default()
	r.GET("/", index)
	r.POST("/:id", processHandle)
	r.GET("/:id", processHandle)
	apiRouter := r.Group("/api")
	apiRouter.POST("/create", api.create)
	apiRouter.POST("/update", api.update)

	r.Static("./static", "./static")
	r.LoadHTMLGlob("./home/*")
	r.Run(":80") //监听端口

}
