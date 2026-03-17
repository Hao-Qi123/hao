package service

import (
	"fmt"
	"seckil/mq/publish"

	"github.com/gin-gonic/gin"
)

func NotifyPay(c *gin.Context) {
	c.Request.ParseForm()
	TradeStatus := c.Request.PostForm.Get("trade_status")
	if TradeStatus != "TRADE_SUCCESS" {
		return
	}
	outTradeNo := c.Request.PostForm.Get("out_trade_no")
	if outTradeNo == "" {
		return
	}
	go func() {
		publish.SendMessage("topic", outTradeNo)
		fmt.Println(outTradeNo)
	}()
	c.JSON(200, gin.H{
		"msg": "success",
	})
}
