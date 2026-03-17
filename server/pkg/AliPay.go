package pkg

import (
	"fmt"
	"seckil/server/basic/config"
	"strconv"

	"github.com/smartwalle/alipay/v3"
)

func Alipay(orderNo string, totalAmount float64) string {
	var configAlipay = config.GlobalConfig.AliPay
	var privateKey = configAlipay.PrivateKey // 必须，上一步中使用 RSA签名验签工具 生成的私钥
	var client, err = alipay.New(configAlipay.AppID, privateKey, false)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	var p = alipay.TradePagePay{}
	p.NotifyURL = configAlipay.NotifyUrl                         // 异步通知URL（支付成功后，支付平台回调的地址）
	p.ReturnURL = configAlipay.ReturnUrl                         // 同步返回URL（支付完成后，页面跳转返回的地址）
	p.Subject = "支付宝支付"                                          // 订单标题/商品名称
	p.OutTradeNo = orderNo                                       // 商户订单号（唯一标识）
	p.TotalAmount = strconv.FormatFloat(totalAmount, 'f', 2, 64) // 订单金额（格式化为保留两位小数的字符串）
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"                     // 产品代码（固定值，表示手机网站支付）

	var url, _ = client.TradePagePay(p)
	// 这个 payURL 即是用于打开支付宝支付页面的 URL，可将输出的内容复制，到浏览器中访问该 URL 即可打开支付页面。
	var payURL = url.String()
	fmt.Println(payURL)
	return payURL
}
