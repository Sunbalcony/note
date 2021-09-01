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
	"time"
)

var DBEngine *gorm.DB

type Message struct {
	Tag  string `json:"tag"`
	Text string `json:"text"`
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
	ctx.Redirect(http.StatusMovedPermanently, ctx.Request.URL.String()+RandChar(8)) //url path长度

}

func processHandle(ctx *gin.Context) {
	msg := &Message{}
	if ctx.Request.Method == http.MethodPost {
		fmt.Println(ctx.Request.FormValue("text"))
		//dataSet[ctx.Request.RequestURI] = ctx.PostForm("text")
		if DBEngine.Where("tag = ?", ctx.Request.URL.Path).First(&msg).RecordNotFound() == true {
			msg.Tag = ctx.Request.URL.Path
			msg.Text = ctx.Request.FormValue("text")
			DBEngine.Create(&msg)
		}
		if DBEngine.Where("tag = ?", ctx.Request.URL.Path).First(&msg).RecordNotFound() == false {
			DBEngine.Model(&msg).Where("tag = ?", ctx.Request.URL.Path).Update("text", ctx.Request.FormValue("text"))
		}

	}
	if ctx.Request.Method == http.MethodGet {
		if ctx.Request.URL.Path != "/favicon.ico" {
			DBEngine.Where("tag = ?", ctx.Request.URL.Path).First(&msg)
			ctx.HTML(http.StatusOK, "index.html", gin.H{"title": "免登录 可分享 实时保存的记事本", "text": msg.Text})
		}

	}

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
		"dbusername", //用户名
		"dbpassword", //密码
		"dburl:3306", //db地址
		"notes",      //库名，要先建库
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
	r := gin.Default()
	r.GET("/", index)
	r.Any("/:id", processHandle)
	r.Static("./static", "./static")
	r.LoadHTMLGlob("./home/*")
	r.Run(":23456")

}
