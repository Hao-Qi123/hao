// [新增文件] service/stock.go
package service

import (
	"fmt"
	"log"
	"sync"
)

// 库存服务
type StockService struct {
	mu     sync.Mutex
	stocks map[string]int
}

// 扣减请求
type ReduceRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	OrderID   string `json:"order_id"`
}

func NewStockService() *StockService {
	return &StockService{
		stocks: make(map[string]int),
	}
}

// 初始化库存
func (s *StockService) InitStock(productID string, num int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stocks[productID] = num
	log.Printf("初始化库存 %s=%d", productID, num)
}

// 扣减库存
func (s *StockService) ReduceStock(req ReduceRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, ok := s.stocks[req.ProductID]
	if !ok {
		return fmt.Errorf("商品不存在")
	}
	if current < req.Quantity {
		return fmt.Errorf("库存不足")
	}

	s.stocks[req.ProductID] = current - req.Quantity

	// [新增] 记录日志
	log.Printf("扣减库存成功 订单:%s 商品:%s 数量:%d 剩余:%d",
		req.OrderID, req.ProductID, req.Quantity, s.stocks[req.ProductID])

	return nil
}

func (s *StockService) GetStock(pid string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.stocks[pid]
}
