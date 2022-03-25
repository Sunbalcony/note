package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
	"net/url"
)

var DBEngine *gorm.DB

func NewDBEngine() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=%s",
		viper.GetString("note.mysqlUsername"),     //用户名
		viper.GetString("note.mysqlPassword"),     //密码
		viper.GetString("note.mysqlUrl"),          //db地址
		viper.GetString("note.mysqlDatabasename"), //库名，要先建库
		"utf8",
		true,
		url.QueryEscape(viper.GetString("note.timezone")),
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
