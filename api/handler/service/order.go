package service

import (
	"net/http"
	"seckil/api/basic/config"
	"seckil/api/handler/request"
	__ "seckil/proto"

	"github.com/gin-gonic/gin"
)

func CreateOrder(c *gin.Context) {
	var req request.OrderReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	r, err := config.ProductClient.AddOrders(c, &__.AddOrdersReq{
		Id:            req.Id,
		OrderNo:       req.OrderNo,
		MemberId:      req.MemberId,
		ProductId:     req.ProductId,
		Consignee:     req.Consignee,
		Mobile:        req.Mobile,
		Address:       req.Address,
		TotalAmount:   req.TotalAmount,
		PayAmount:     req.PayAmount,
		FreightAmount: req.FreightAmount,
		PaymentMethod: req.PaymentMethod,
		PaymentAt:     req.PaymentAt,
		Status:        req.Status,
		OrderType:     req.OrderType,
		Remark:        req.Remark,
		CancelReason:  req.CancelReason,
		Quantity:      req.Quantity,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ID":          r.ID,
		"OrderNo":     r.OrderNo,
		"ProductID":   r.ProductID,
		"ProductName": r.ProductName,
		"Quantity":    r.Quantity,
		"TotalAmount": r.TotalAmount,
		"Status":      r.Status,
		"Remark":      r.Remark,
		"PayUrl":      r.PayUrl,
		"Msg":         "订单创建成功",
	})
}
