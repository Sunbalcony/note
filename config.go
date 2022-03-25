package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
)

func InitConfig() {
	workdir, _ := os.Getwd()
	fmt.Println("当前目录", workdir)
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workdir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

}
