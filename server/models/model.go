package models

import (
	"time"

	"gorm.io/gorm"
)

//会员表 (member)

type Member struct {
	ID          uint       `gorm:"primaryKey;comment:'会员ID'"`
	Username    string     `gorm:"type:varchar(50);uniqueIndex;not null;comment:'用户名'"`
	Password    string     `gorm:"type:varchar(100);not null;comment:'密码'"`
	Mobile      string     `gorm:"type:varchar(20);uniqueIndex;comment:'手机号'"`
	Email       string     `gorm:"type:varchar(100);comment:'邮箱'"`
	Avatar      string     `gorm:"type:varchar(255);comment:'头像URL'"`
	Nickname    string     `gorm:"type:varchar(50);comment:'昵称'"`
	RealName    string     `gorm:"type:varchar(50);comment:'真实姓名'"`
	Gender      int        `gorm:"type:tinyint(1);default:0;comment:'性别:0未知1男2女'"`
	Birthday    *time.Time `gorm:"type:date;comment:'生日'"`
	Status      int        `gorm:"type:tinyint(1);default:1;comment:'状态:0禁用1启用'"`
	LastLoginIP string     `gorm:"type:varchar(50);comment:'最后登录IP'"`
	LastLoginAt *time.Time `gorm:"type:datetime;comment:'最后登录时间'"`
	CreatedAt   time.Time  `gorm:"comment:'注册时间'"`
	UpdatedAt   time.Time  `gorm:"comment:'更新时间'"`
}

//商品表 (product)

type Product struct {
	ID          uint      `gorm:"primaryKey;comment:'商品ID'"`
	CategoryID  uint      `gorm:"index;comment:'分类ID'"`
	ProductNo   string    `gorm:"type:varchar(50);uniqueIndex;not null;comment:'商品编号'"`
	Name        string    `gorm:"type:varchar(200);not null;index;comment:'商品名称'"`
	Keywords    string    `gorm:"type:varchar(255);comment:'关键词'"`
	Description string    `gorm:"type:text;comment:'商品描述'"`
	MainImage   string    `gorm:"type:varchar(500);comment:'主图'"`
	Price       float64   `gorm:"type:decimal(10,2);not null;comment:'价格'"`
	SalesCount  int       `gorm:"default:0;comment:'销量'"`
	ReviewCount int       `gorm:"default:0;comment:'评价数'"`
	Status      int       `gorm:"type:tinyint(1);default:1;comment:'状态:0下架1上架'"`
	Stock       int       `gorm:"default:200;comment:'库存'"`
	Quantity    int       `gorm:"not null;comment:'数量'"`
	CreatedAt   time.Time `gorm:"comment:'创建时间'"`
	UpdatedAt   time.Time `gorm:"comment:'更新时间'"`
}

func (p *Product) CreateProduct(db *gorm.DB) error {
	return db.Debug().Create(&p).Error
}

func (p *Product) UpdateProduct(db *gorm.DB, id int64, m map[string]interface{}) error {
	return db.Debug().Model(&p).Where("id=?", id).Updates(&m).Error
}

func (p *Product) DeleteProduct(db *gorm.DB, id int64) error {
	return db.Debug().Where("id=?", id).Delete(&p).Error
}

func (p *Product) GetProductById(db *gorm.DB, id int64) error {
	return db.Debug().Where("id=?", id).First(&p).Error
}

func (p *Product) FirstProduct(db *gorm.DB, id int64) error {
	return db.Debug().Where("id=?", id).First(&p).Error
}

//库存表 (stock)

type Stock struct {
	ID             uint       `gorm:"primaryKey;comment:'库存ID'"`
	ProductID      uint       `gorm:"uniqueIndex;not null;comment:'商品ID'"`
	SkuID          uint       `gorm:"index;comment:'SKU ID'"`
	WarehouseID    uint       `gorm:"comment:'仓库ID'"`
	TotalStock     int        `gorm:"default:0;comment:'总库存'"`
	AvailableStock int        `gorm:"default:0;comment:'可用库存'"`
	FrozenStock    int        `gorm:"default:0;comment:'冻结库存(已下单未支付)'"`
	LockedStock    int        `gorm:"default:0;comment:'锁定库存(已支付)'"`
	AlertStock     int        `gorm:"default:0;comment:'预警库存'"`
	MaxStock       int        `gorm:"default:0;comment:'最大库存'"`
	MinStock       int        `gorm:"default:0;comment:'最小库存'"`
	LastInTime     *time.Time `gorm:"type:datetime;comment:'最后入库时间'"`
	LastOutTime    *time.Time `gorm:"type:datetime;comment:'最后出库时间'"`
	CreatedAt      time.Time  `gorm:"comment:'创建时间'"`
	UpdatedAt      time.Time  `gorm:"comment:'更新时间'"`
}

//积分表 (point)

type Point struct {
	ID        uint       `gorm:"primaryKey;comment:'积分ID'"`
	MemberID  uint       `gorm:"index;not null;comment:'会员ID'"`
	OrderID   uint       `gorm:"index;comment:'订单ID'"`
	PointNo   string     `gorm:"type:varchar(50);uniqueIndex;comment:'积分流水号'"`
	Type      int        `gorm:"type:tinyint(1);comment:'类型:1获取2消费3过期4退款'"`
	Point     int        `gorm:"not null;comment:'积分数值'"`
	Balance   int        `gorm:"default:0;comment:'积分余额'"`
	Source    string     `gorm:"type:varchar(50);comment:'来源:购物/签到/活动'"`
	ExpireAt  *time.Time `gorm:"type:datetime;comment:'过期时间'"`
	Remark    string     `gorm:"type:varchar(255);comment:'备注'"`
	CreatedAt time.Time  `gorm:"comment:'创建时间'"`
}

//订单表 (order)

type Order struct {
	ID            uint       `gorm:"primaryKey;comment:'订单ID'"`
	OrderNo       string     `gorm:"type:varchar(50);uniqueIndex;not null;comment:'订单号'"`
	MemberID      uint       `gorm:"index;not null;comment:'会员ID'"`
	ProductId     int64      `gorm:"type:int(10);not null;comment:'商品ID'"`
	Consignee     string     `gorm:"type:varchar(50);not null;comment:'收货人'"`
	Mobile        string     `gorm:"type:varchar(20);not null;comment:'联系电话'"`
	Address       string     `gorm:"type:varchar(255);not null;comment:'详细地址'"`
	TotalAmount   float64    `gorm:"type:decimal(10,2);not null;comment:'订单总额'"`
	PayAmount     float64    `gorm:"type:decimal(10,2);not null;comment:'实付金额'"`
	FreightAmount float64    `gorm:"type:decimal(10,2);default:0;comment:'运费'"`
	PaymentMethod int        `gorm:"type:tinyint(1);comment:'支付方式:1微信2支付宝3积分'"`
	PaymentAt     *time.Time `gorm:"type:datetime;comment:'支付时间'"`
	Status        int        `gorm:"type:tinyint(2);default:1;comment:'状态:1待付款2待发货3待收货4已完成5已取消6售后'"`
	OrderType     int        `gorm:"type:tinyint(1);default:1;comment:'订单类型:1普通2秒杀3团购'"`
	Remark        string     `gorm:"type:varchar(500);comment:'订单备注'"`
	CancelReason  string     `gorm:"type:varchar(255);comment:'取消原因'"`
	CreatedAt     time.Time  `gorm:"comment:'下单时间'"`
	UpdatedAt     time.Time  `gorm:"comment:'更新时间'"`
}

//订单商品表 (order_item)

type OrderItem struct {
	ID             uint      `gorm:"primaryKey;comment:'订单商品ID'"`
	OrderID        uint      `gorm:"index;not null;comment:'订单ID'"`
	ProductID      uint      `gorm:"not null;comment:'商品ID'"`
	ProductName    string    `gorm:"type:varchar(200);not null;comment:'商品名称'"`
	Image          string    `gorm:"type:varchar(500);comment:'商品图片'"`
	Price          float64   `gorm:"type:decimal(10,2);not null;comment:'单价'"`
	Quantity       int       `gorm:"not null;comment:'数量'"`
	TotalAmount    float64   `gorm:"type:decimal(10,2);not null;comment:'总价'"`
	DiscountAmount float64   `gorm:"type:decimal(10,2);default:0;comment:'优惠金额'"`
	PayAmount      float64   `gorm:"type:decimal(10,2);not null;comment:'实付金额'"`
	Point          int       `gorm:"default:0;comment:'赠送积分'"`
	IsReviewed     int       `gorm:"type:tinyint(1);default:0;comment:'是否已评价'"`
	CreatedAt      time.Time `gorm:"comment:'创建时间'"`
}

//物流配送表 (logistics)

type Logistics struct {
	ID              uint       `gorm:"primaryKey;comment:'物流ID'"`
	OrderID         uint       `gorm:"uniqueIndex;not null;comment:'订单ID'"`
	LogisticsNo     string     `gorm:"type:varchar(50);uniqueIndex;comment:'物流单号'"`
	CompanyCode     string     `gorm:"type:varchar(50);comment:'快递公司编码'"`
	CompanyName     string     `gorm:"type:varchar(100);comment:'快递公司名称'"`
	ShipperName     string     `gorm:"type:varchar(50);comment:'发货人'"`
	ShipperMobile   string     `gorm:"type:varchar(20);comment:'发货人电话'"`
	ShipperAddress  string     `gorm:"type:varchar(255);comment:'发货地址'"`
	ReceiverName    string     `gorm:"type:varchar(50);comment:'收货人'"`
	ReceiverMobile  string     `gorm:"type:varchar(20);comment:'收货人电话'"`
	ReceiverAddress string     `gorm:"type:varchar(255);comment:'收货地址'"`
	Weight          float64    `gorm:"type:decimal(10,2);comment:'包裹重量'"`
	Volume          float64    `gorm:"type:decimal(10,2);comment:'包裹体积'"`
	Fee             float64    `gorm:"type:decimal(10,2);comment:'运费'"`
	Status          int        `gorm:"type:tinyint(1);default:0;comment:'状态:0待发货1已发货2配送中3已签收4异常'"`
	TrackingInfo    string     `gorm:"type:text;comment:'物流跟踪信息(JSON)'"`
	ShipAt          *time.Time `gorm:"type:datetime;comment:'发货时间'"`
	ReceiveAt       *time.Time `gorm:"type:datetime;comment:'签收时间'"`
	CreatedAt       time.Time  `gorm:"comment:'创建时间'"`
	UpdatedAt       time.Time  `gorm:"comment:'更新时间'"`
}
