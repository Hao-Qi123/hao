package service

import (
	"net/http"
	"seckil/api/basic/config"
	"seckil/api/handler/request"
	__ "seckil/proto"

	"github.com/gin-gonic/gin"
)

func CreateProduct(c *gin.Context) {
	var req request.AddProductReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	_, err = config.ProductClient.AddProducts(c, &__.AddProductsReq{
		ProductNo:   req.ProductNo,
		Name:        req.Name,
		Keywords:    req.Keywords,
		Description: req.Description,
		MainImage:   req.MainImage,
		Price:       req.Price,
		SalesCount:  req.SalesCount,
		ReviewCount: req.ReviewCount,
		Status:      req.Status,
		Stock:       req.Stock,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建成功",
	})
}
func UpdateProduct(c *gin.Context) {
	var req request.UpdateProductReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	_, err = config.ProductClient.UpdateProducts(c, &__.UpdateProductsReq{
		ProductNo:   req.ProductNo,
		Name:        req.Name,
		Keywords:    req.Keywords,
		Description: req.Description,
		MainImage:   req.MainImage,
		Price:       req.Price,
		SalesCount:  req.SalesCount,
		ReviewCount: req.ReviewCount,
		Status:      req.Status,
		Stock:       req.Stock,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "修改成功",
	})

}

func DeleteProduct(c *gin.Context) {
	var req request.DelProductReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	_, err = config.ProductClient.DelProducts(c, &__.DelProductsReq{
		Id: req.ID,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})

}
func GetProductByID(c *gin.Context) {
	var req request.GetProductReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	r, err := config.ProductClient.GetProductsById(c, &__.GetProductsByIdReq{
		Id: req.ID,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "查询成功",
		"product": r.Products,
	})
}
func ListProducts(c *gin.Context) {
	var req request.FindProductReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	r, err := config.ProductClient.SearchProducts(c, &__.SearchProductsReq{
		Page:      req.Page,
		Size:      req.Size,
		ProductNo: req.ProductNo,
		Name:      req.Name,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"msg":      "查询成功",
		"products": r.Products,
	})
}
