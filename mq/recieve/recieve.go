package main

import (
	"context"
	"fmt"
	"log"
	config2 "seckil/api/basic/config"
	"seckil/mq/rabbitMQ"
	__ "seckil/proto"
	_ "seckil/server/basic/bootstrap"
	"seckil/server/basic/config"
	"time"
)

func main() {

	kutengOne := rabbitMQ.NewRabbitMQTopic("exKutengTopic", "kuteng.topic.one")
	kutengOne.RecieveTopic("topic", func(msg string) {
		result, _ := config.Rdb.SetNX(config.Ctx, "order_sn", 1, time.Minute*60).Result()
		if !result {
			log.Println("1111")
			fmt.Println(msg)
			return
		}

		fmt.Println(msg)
		config2.ProductClient.UpdatesStatus(context.Background(), &__.UpdatesStatusReq{
			OrderSn: msg,
		})

	})
}
