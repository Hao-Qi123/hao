package request

type AddProductReq struct {
	CategoryId  int64   `form:"categoryId" json:"categoryId" xml:"categoryId" binding:"required"`
	ProductNo   string  `form:"productNo" json:"productNo" xml:"productNo" binding:"required"`
	Name        string  `form:"name" json:"name" xml:"name" binding:"required"`
	Keywords    string  `form:"keywords" json:"keywords" xml:"keywords" binding:"required"`
	Description string  `form:"description" json:"description" xml:"description" binding:"required"`
	MainImage   string  `form:"mainImage" json:"mainImage" xml:"mainImage" binding:"required"`
	Price       float64 `form:"price" json:"price" xml:"price" binding:"required"`
	SalesCount  int64   `form:"salesCount" json:"salesCount" xml:"salesCount" binding:"required"`
	ReviewCount int64   `form:"reviewCount" json:"reviewCount" xml:"reviewCount" binding:"required"`
	Status      int64   `form:"status" json:"status" xml:"status" binding:"required"`
	Stock       int64   `form:"stock" json:"stock" xml:"stock" binding:"required"`
}
type UpdateProductReq struct {
	ID          int64   `form:"id" json:"id" xml:"id" binding:"required"`
	CategoryId  int64   `form:"categoryId" json:"categoryId" xml:"categoryId" binding:"required"`
	ProductNo   string  `form:"productNo" json:"productNo" xml:"productNo" binding:"required"`
	Name        string  `form:"name" json:"name" xml:"name" binding:"required"`
	Keywords    string  `form:"keywords" json:"keywords" xml:"keywords" binding:"required"`
	Description string  `form:"description" json:"description" xml:"description" binding:"required"`
	MainImage   string  `form:"mainImage" json:"mainImage" xml:"mainImage" binding:"required"`
	Price       float64 `form:"price" json:"price" xml:"price" binding:"required"`
	SalesCount  int64   `form:"salesCount" json:"salesCount" xml:"salesCount" binding:"required"`
	ReviewCount int64   `form:"reviewCount" json:"reviewCount" xml:"reviewCount" binding:"required"`
	Status      int64   `form:"status" json:"status" xml:"status" binding:"required"`
	Stock       int64   `form:"stock" json:"stock" xml:"stock" binding:"required"`
}
type DelProductReq struct {
	ID int64 `form:"id" json:"id" xml:"id" binding:"required"`
}
type GetProductReq struct {
	ID int64 `form:"id" json:"id" xml:"id" binding:"required"`
}
type FindProductReq struct {
	Page      int64  `form:"page" json:"page" xml:"page"`
	Size      int64  `form:"size" json:"size" xml:"size"`
	ProductNo string `form:"productNo" json:"productNo" xml:"productNo"`
	Name      string `form:"name" json:"name" xml:"name"`
}
type OrderReq struct {
	Id            int64   `form:"id" json:"id" xml:"id" binding:"required"`
	OrderNo       string  `form:"orderNo" json:"orderNo" xml:"orderNo" binding:""`
	MemberId      int64   `form:"memberId" json:"memberId" xml:"memberId" binding:"required"`
	ProductId     int64   `form:"productId" json:"productId" xml:"productId" binding:"required"`
	Consignee     string  `form:"consignee" json:"consignee" xml:"consignee" binding:"required"`
	Mobile        string  `form:"mobile" json:"mobile" xml:"mobile" binding:"required"`
	Address       string  `form:"address" json:"address" xml:"address" binding:"required"`
	TotalAmount   float64 `form:"totalAmount" json:"totalAmount" xml:"totalAmount" binding:"required"`
	PayAmount     float64 `form:"payAmount" json:"payAmount" xml:"payAmount" binding:"required"`
	FreightAmount float64 `form:"freightAmount" json:"freightAmount" xml:"freightAmount" binding:"required"`
	PaymentMethod int64   `form:"paymentMethod" json:"paymentMethod" xml:"paymentMethod" binding:"required"`
	PaymentAt     int64   `form:"paymentAt" json:"paymentAt" xml:"paymentAt" binding:"required"`
	Status        int64   `form:"status" json:"status" xml:"status" binding:"required"`
	OrderType     int64   `form:"orderType" json:"orderType" xml:"orderType" binding:"required"`
	Remark        string  `form:"remark" json:"remark" xml:"remark" binding:"required"`
	CancelReason  string  `form:"cancelReason" json:"cancelReason" xml:"cancelReason" binding:"required"`
	Quantity      int64   `form:"quantity" json:"quantity" xml:"quantity" binding:"required"`
}
