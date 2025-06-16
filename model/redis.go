package model

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"math/rand"
	"time"
)

type Redis struct {
	rds *redis.Client
}

var Rds *Redis

func NewRedisApi() *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:        viper.GetString("note.redisUrl"),      // redis地址
		Password:    viper.GetString("note.redisPassword"), // redis密码，没有则留空
		DB:          viper.GetInt("note.redisDatabaseNum"), // 默认数据库，默认是0
		PoolSize:    100000,                                // Redis连接池大小
		MaxRetries:  2,                                     // 最大重试次数
		IdleTimeout: 10 * time.Second,                      // 空闲链接超时时间
	})
	return &Redis{rds: client}

}

func (r *Redis) appendKey(oneKey string, oneValue string) {
	err := r.rds.Append(oneKey, oneValue).Err()
	// 检测设置是否成功
	if err != nil {
		fmt.Printf("append出错%s", err)
		panic(err)
	}

}

func (r *Redis) SetKey(oneKey string, oneValue string) error {
	// 设置一个key，过期时间为0，意思就是永远不过期
	err := r.rds.Set(oneKey, oneValue, 0).Err()

	// 检测设置是否成功
	if err != nil {
		fmt.Printf("set出错%s", err)
		return err
	}
	return nil

}
func (r *Redis) HSetKey(oneKey string, oneField string, oneValue string) {
	err := r.rds.HSet(oneKey, oneField, oneValue).Err()
	if err != nil {
		fmt.Printf("Hset:%s域出错%s", oneField, err)
	}

}
func (r *Redis) HSetKeyList(oneKey string, oneField string, oneValue []string) {

	aaa, _ := json.Marshal(oneValue)
	fmt.Printf("A是%s", aaa)

	err := r.rds.HSet(oneKey, oneField, aaa).Err()
	if err != nil {
		fmt.Printf("Hset:%s域出错%s", oneField, err)
	}

}

func (r *Redis) GetKey(oneExistKey string) (string, error) {
	// 根据key查询缓存，通过Result函数返回两个值
	//  第一个代表key的值，第二个代表查询错误信息
	val, err := r.rds.Get(oneExistKey).Result()

	// 检测，查询是否出错
	if err != nil {
		return "", err
	}
	return val, nil

}
func (r *Redis) HGetKey(oneKey string, oneField string) string {
	val, err := r.rds.HGet(oneKey, oneField).Result()
	if err != nil {
		fmt.Printf("HGet%s;%s值异常", oneKey, oneField)
	}
	fmt.Println("HGet的值是", val)
	return val

}
func (r *Redis) KeyExist(oneKey string) int64 {
	//存在返回1，不存在返回0
	val, err := r.rds.Exists(oneKey).Result()
	if err != nil {
		fmt.Println(err)
	}
	return val

}
func (r *Redis) HKeyExist(oneKey, oneField string) bool {
	val, err := r.rds.HExists(oneKey, oneField).Result()
	if err != nil {
		fmt.Println(err)

	}
	return val

}
func (r *Redis) HKeyValues(oneKey string) []string {
	val, err := r.rds.HVals(oneKey).Result()
	if err != nil {
		fmt.Println(err)

	}
	fmt.Println(val)
	return val

}
func (r *Redis) getAllKeys() []string {
	val, err := r.rds.Keys("*").Result()
	if err != nil {
		fmt.Println(err)

	}
	fmt.Println(val)
	return val

}
func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	fmt.Printf("全局标识参数:%s\n", string(result))
	return string(result)
}
