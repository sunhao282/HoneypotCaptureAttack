package conf

import (
	"fmt"
	"gopkg.in/ini.v1"
)

var cfg *ini.File

//#初始化函数
func init() {
	//读取ini配置文件
	c, err := ini.Load("config.ini")
	if err != nil {
		fmt.Println("打开配置文件失败")
	}
	c.BlockMode = false
	cfg = c
}

func Get(node string, key string) string {
	val := cfg.Section(node).Key(key).String()
	fmt.Println("val:", val)
	return val
}

func GetInt(node string, key string) int {
	val, _ := cfg.Section(node).Key(key).Int()
	return val
}
