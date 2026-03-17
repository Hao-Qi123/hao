package bootstrap

import (
	"fmt"
	"seckil/server/basic/config"

	"github.com/spf13/viper"
)

func InitAppConfig() {
	viper.SetConfigFile("E:\\go\\gocode\\src\\seckill\\config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic("viper配置读取失败" + err.Error())
	}
	err = viper.Unmarshal(&config.GlobalConfig)
	if err != nil {
		panic("viper配置解析失败" + err.Error())
	}
	fmt.Println("viper配置成功", config.GlobalConfig)
}
