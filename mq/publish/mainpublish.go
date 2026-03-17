package publish

import (
	"seckil/mq/rabbitMQ"
)

func SendMessage(topic string, msg string) {
	kutengOne := rabbitMQ.NewRabbitMQTopic("exKutengTopic", "kuteng.topic.one")
	kutengOne.PublishTopic(topic, msg)

}
