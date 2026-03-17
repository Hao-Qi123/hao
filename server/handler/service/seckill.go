package service

import (
	"context"
	"log"
	"seckil/server/basic/config"
	"seckil/server/models"
	"seckil/server/pkg"
	"time"

	__ "seckil/proto"

	"gorm.io/gorm"
)

// server is used to implement helloworld.GreeterServer.
type Server struct {
	__.UnimplementedProductServer
}

// SayHello implements helloworld.GreeterServer
func (s *Server) AddProducts(_ context.Context, in *__.AddProductsReq) (*__.AddProductsResp, error) {
	log.Printf("Received: %v", in.GetName())
	product := models.Product{
		CategoryID:  uint(in.CategoryId),
		ProductNo:   pkg.OrderNo(),
		Name:        in.Name,
		Keywords:    in.Keywords,
		Description: in.Description,
		MainImage:   in.MainImage,
		Price:       in.Price,
		SalesCount:  int(in.SalesCount),
		ReviewCount: int(in.ReviewCount),
		Status:      int(in.Status),
		Stock:       int(in.Stock),
		Quantity:    int(in.Quantity),
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	}
	err := product.CreateProduct(config.DB)
	if err != nil {
		return &__.AddProductsResp{
			Code: 400,
			Msg:  "创建商品失败",
		}, nil
	}
	return &__.AddProductsResp{
		Code: 200,
		Msg:  "创建商品成功",
	}, nil
}

// SayHello implements helloworld.GreeterServer
func (s *Server) UpdateProducts(_ context.Context, in *__.UpdateProductsReq) (*__.UpdateProductsResp, error) {
	log.Printf("Received: %v", in.GetName())
	var product models.Product
	product.ID = uint(in.Id)
	err := product.UpdateProduct(config.DB, in.Id, map[string]interface{}{
		"ID":          in.Id,
		"CategoryID":  uint(in.CategoryId),
		"ProductNo":   in.ProductNo,
		"Name":        in.Name,
		"Keywords":    in.Keywords,
		"Description": in.Description,
		"MainImage":   in.MainImage,
		"Price":       in.Price,
		"SalesCount":  int(in.SalesCount),
		"ReviewCount": int(in.ReviewCount),
		"Status":      int(in.Status),
		"Stock":       int(in.Stock),
		"Quantity":    int(in.Quantity),
	})
	if err != nil {
		return &__.UpdateProductsResp{
			Code: 400,
			Msg:  "更新商品失败",
		}, nil
	}
	return &__.UpdateProductsResp{
		Code: 200,
		Msg:  "更新商品成功",
	}, nil
}

// SayHello implements helloworld.GreeterServer
func (s *Server) DelProducts(_ context.Context, in *__.DelProductsReq) (*__.DelProductsResp, error) {
	var product models.Product
	err := product.DeleteProduct(config.DB, in.Id)
	if err != nil {
		return &__.DelProductsResp{
			Code: 400,
			Msg:  "删除商品失败",
		}, nil
	}
	return &__.DelProductsResp{
		Code: 200,
		Msg:  "删除商品成功",
	}, nil
}

// SayHello implements helloworld.GreeterServer
func (s *Server) GetProductsById(_ context.Context, in *__.GetProductsByIdReq) (*__.GetProductsByIdResp, error) {
	var product models.Product
	err := product.GetProductById(config.DB, in.Id)
	if err != nil {
		return &__.GetProductsByIdResp{
			Code: 400,
			Msg:  "查询商品失败",
		}, nil
	}
	var createdAt, updatedAt int64
	if !product.CreatedAt.IsZero() {
		createdAt = product.CreatedAt.Unix()
	}
	if !product.UpdatedAt.IsZero() {
		updatedAt = product.UpdatedAt.Unix()
	}

	return &__.GetProductsByIdResp{
		Code: 200,
		Msg:  "查询商品成功",
		Products: &__.Products{
			Id:          int64(product.ID),
			CategoryId:  int64(product.CategoryID),
			ProductNo:   product.ProductNo,
			Name:        product.Name,
			Keywords:    product.Keywords,
			Description: product.Description,
			MainImage:   product.MainImage,
			Price:       product.Price,
			SalesCount:  int64(product.SalesCount),
			ReviewCount: int64(product.ReviewCount),
			Status:      int64(product.Status),
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		},
	}, nil
}

// SayHello implements helloworld.GreeterServer
func (s *Server) SearchProducts(_ context.Context, in *__.SearchProductsReq) (*__.SearchProductsResp, error) {
	// 1. 参数验证和默认值设置
	page := in.GetPage()
	size := in.GetSize()

	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100 // 限制最大查询数量
	}

	offset := (page - 1) * size

	// 2. 初始化查询
	var products []models.Product

	db := config.DB.Model(&models.Product{})

	// 3. 构建查询条件
	if productNo := in.GetProductNo(); productNo != "" {
		db = db.Where("product_no LIKE ?", "%"+productNo+"%")
	}

	if name := in.GetName(); name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}

	// 5. 执行分页查询
	if err := db.Offset(int(offset)).Limit(int(size)).Find(&products).Error; err != nil {
		return &__.SearchProductsResp{
			Code: 500,
			Msg:  "查询商品失败",
		}, nil
	}

	// 6. 转换为响应格式
	var productsResp []*__.Products
	for _, product := range products {
		// 处理时间字段
		createdAt := product.CreatedAt.Unix()
		updatedAt := product.UpdatedAt.Unix()

		productsResp = append(productsResp, &__.Products{
			Id:          int64(product.ID),
			CategoryId:  int64(product.CategoryID),
			ProductNo:   product.ProductNo,
			Name:        product.Name,
			Keywords:    product.Keywords,
			Description: product.Description,
			MainImage:   product.MainImage,
			Price:       product.Price,
			SalesCount:  int64(product.SalesCount),
			ReviewCount: int64(product.ReviewCount),
			Status:      int64(product.Status),
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		})
	}

	// 7. 返回结果（包含分页信息）
	return &__.SearchProductsResp{
		Code:     200,
		Msg:      "查询成功",
		Products: productsResp,
		Page:     page, // 可能需要添加 page 字段
		Size:     size, // 可能需要添加 size 字段
	}, nil
}
func (s *Server) AddOrders(_ context.Context, in *__.AddOrdersReq) (*__.AddOrdersResp, error) {
	orderNo := pkg.OrderNo()
	var product models.Product

	// 1. 先查询商品信息
	err := product.FirstProduct(config.DB, in.ProductId)
	if err != nil {
		return &__.AddOrdersResp{
			Msg: "商品不存在",
		}, nil
	}
	// 2. 检查库存
	if product.Stock < int(in.Quantity) {
		return &__.AddOrdersResp{
			Msg: "商品库存不足",
		}, nil
	}

	// 3. 计算订单金额
	totalAmount := product.Price * float64(in.Quantity)
	payAmount := totalAmount + in.FreightAmount - in.FreightAmount

	// 4. 创建主订单
	order := models.Order{
		OrderNo:       orderNo,
		MemberID:      uint(in.MemberId),
		ProductId:     in.ProductId,
		Consignee:     in.Consignee,
		Mobile:        in.Mobile,
		Address:       in.Address,
		TotalAmount:   totalAmount,
		PayAmount:     payAmount,
		FreightAmount: in.FreightAmount,
		PaymentMethod: int(in.PaymentMethod),
		PaymentAt:     nil,
		Status:        int(in.Status), // 1-待付款
		OrderType:     int(in.OrderType),
		Remark:        in.Remark,
		CancelReason:  in.CancelReason,
		CreatedAt:     time.Time{},
		UpdatedAt:     time.Time{},
	}

	// 5. 开启事务
	tx := config.DB.Begin()

	// 6. 创建订单
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return &__.AddOrdersResp{
			Msg: "创建订单失败",
		}, nil
	}

	// 7. 创建订单商品记录
	orderItem := models.OrderItem{
		OrderID:        order.ID,
		ProductID:      product.ID,
		ProductName:    product.Name,
		Image:          product.MainImage,
		Price:          product.Price,
		Quantity:       int(in.Quantity),
		TotalAmount:    totalAmount,
		DiscountAmount: 3,
		PayAmount:      payAmount,
		Point:          2,
		IsReviewed:     1,
	}

	if err := tx.Create(&orderItem).Error; err != nil {
		tx.Rollback()
		return &__.AddOrdersResp{
			Msg: "创建订单商品失败",
		}, nil
	}

	// 8. 扣减库存
	if err := tx.Model(&models.Product{}).Where("id = ? AND stock >= ?", in.ProductId, in.Quantity).
		Update("stock", gorm.Expr("stock - ?", in.Quantity)).Error; err != nil {
		tx.Rollback()
		return &__.AddOrdersResp{
			Msg: "库存不足",
		}, nil
	}

	// 9. 提交事务
	tx.Commit()

	// 10. 生成支付链接
	payUrl := pkg.Alipay(order.OrderNo, order.PayAmount)

	// 11. 返回响应
	return &__.AddOrdersResp{
		ID:          int64(order.ID),
		OrderNo:     order.OrderNo,
		ProductID:   int64(product.ID),
		ProductName: product.Name,
		Quantity:    int64(orderItem.Quantity),
		TotalAmount: float32(order.TotalAmount),
		Status:      int64(order.Status),
		Remark:      order.Remark,
		CreatedAt:   order.CreatedAt.Format("2006-01-02 15:04:05"),
		PayUrl:      payUrl,
		Msg:         "订单创建成功",
	}, nil

}
func (s *Server) UpdatesStatus(_ context.Context, in *__.UpdatesStatusReq) (*__.UpdatesStatusResp, error) {
	// 事务
	tx := config.DB.Begin()
	var order models.Order
	err := tx.Where("order_no = ?", in.OrderSn).First(&order).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	// 幂等
	if order.Status == 2 {
		tx.Commit()

	}
	// 修改订单状态
	order.Status = 2
	err = tx.Save(&order).Error
	if err != nil {
		tx.Rollback()

	}
	// 查询订单明细
	var orderItems []models.OrderItem
	err = tx.Where("id=?", order.ID).Find(&orderItems).Error
	if err != nil {
		tx.Rollback()

	}
	// 扣库存
	for _, items := range orderItems {
		var product models.Product
		err = tx.Where("id=?", items.ProductID).First(&product).Error
		if err != nil {
			tx.Rollback()

		}
		product.Quantity -= items.Quantity
		err = tx.Save(&product).Error
		if err != nil {
			tx.Rollback()

		}
	}
	tx.Commit()
	return &__.UpdatesStatusResp{}, nil
}
