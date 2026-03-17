package bootstrap

import (
	"fmt"
	"seckil/server/basic/config"

	"strings"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
)

func InitNacos() {
	var configNacos config.Nacos
	viper.SetConfigFile("E:\\go\\gocode\\src\\seckill\\config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic("nacos配置读取失败" + err.Error())
	}
	err = viper.UnmarshalKey("nacos", &configNacos)
	if err != nil {
		panic("nacos配置解析失败" + err.Error())
	}
	fmt.Println("nacos配置成功", configNacos)
	// Nacos服务器地址
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: configNacos.Addr,
			Port:   uint64(configNacos.Port),
		},
	}
	// 客户端配置
	clientConfig := constant.ClientConfig{
		NamespaceId:         configNacos.Namespace, // 如果不需要命名空间，可以留空
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "./tmp/nacos/log",
		CacheDir:            "./tmp/nacos/cache",
		LogLevel:            "debug",
	}

	// 创建配置客户端
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	configs, err := configClient.GetConfig(vo.ConfigParam{
		DataId: configNacos.DataId,
		Group:  configNacos.Group,
	})
	err = viper.ReadConfig(strings.NewReader(configs))
	if err != nil {
		return
	}
	fmt.Println(configs)
	err = viper.Unmarshal(&config.GlobalConfig)
	if err != nil {
		return
	}
	fmt.Println("nacos配置成功", config.GlobalConfig)
}
