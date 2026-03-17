package router

import (
	"seckil/api/handler/service"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	product := r.Group("/product")
	{
		product.POST("/add", service.CreateProduct)
		product.POST("/update", service.UpdateProduct)
		product.POST("/delete", service.DeleteProduct)
		product.POST("/get", service.GetProductByID)
		product.POST("/list", service.ListProducts)
	}
	order := r.Group("/order")
	{
		order.POST("/add", service.CreateOrder)
		order.POST("/notify/pay", service.NotifyPay)

	}

	return r
}
