# 秒杀系统业务流程分析文档

> 基于代码库模型定义的业务流程梳理
> 文档版本: v1.0
> 生成时间: 2026-03-14

---

## 目录

1. [订单全流程](#一订单全流程)
2. [库存全流程](#二库存全流程)
3. [积分全流程](#三积分全流程)
4. [数据流程图](#四数据流程图)
5. [核心实体关系](#五核心实体关系)

---

## 一、订单全流程

### 1.1 订单状态定义

基于 `Order` 模型的 `Status` 字段：

| 状态码 | 状态名称 | 说明 |
|--------|----------|------|
| 1 | 待付款 | 订单已创建，等待用户支付 |
| 2 | 待发货 | 已支付，等待商家发货 |
| 3 | 待收货 | 已发货，等待用户确认收货 |
| 4 | 已完成 | 订单完成，交易结束 |
| 5 | 已取消 | 订单被取消（超时/主动取消） |
| 6 | 售后 | 订单进入售后流程 |

### 1.2 订单时序图

```mermaid
sequenceDiagram
    actor U as 用户
    participant C as 客户端/API
    participant O as 订单服务
    participant S as 库存服务
    participant P as 支付服务
    participant L as 物流服务
    participant DB as 数据库

    %% 1. 订单创建
    rect rgb(240, 248, 255)
    Note over U,DB: 阶段1: 订单创建 + 库存预占
    U->>C: 1. 提交订单(商品+数量)
    C->>S: 2. 查询库存可用性
    S->>DB: 查询 Stock.AvailableStock
    DB-->>S: 返回库存数量
    S-->>C: 库存充足
    
    C->>O: 3. 创建订单
    O->>DB: 插入 Order (Status=1)
    O->>DB: 插入 OrderItem
    O->>S: 4. 冻结库存
    S->>DB: Update Stock:<br/>AvailableStock -= qty<br/>FrozenStock += qty
    S-->>O: 冻结成功
    O-->>C: 返回订单信息
    C-->>U: 订单创建成功<br/>(OrderNo, 待付款)
    end

    %% 2. 订单支付
    rect rgb(255, 248, 240)
    Note over U,DB: 阶段2: 订单支付
    U->>C: 5. 发起支付
    C->>P: 6. 创建支付单
    P->>DB: 插入支付记录
    P-->>C: 返回支付链接
    C-->>U: 跳转支付页面
    
    U->>P: 7. 完成支付
    P->>O: 8. 支付回调
    O->>DB: Update Order:<br/>Status=2, PaymentNo, PaymentAt
    O->>S: 9. 库存锁定
    S->>DB: Update Stock:<br/>FrozenStock -= qty<br/>LockedStock += qty
    S-->>O: 锁定成功
    O-->>P: 支付完成确认
    P-->>C: 支付结果通知
    C-->>U: 支付成功
    end

    %% 3. 订单发货
    rect rgb(240, 255, 240)
    Note over U,DB: 阶段3: 订单履约
    U->>C: 10. (可选) 申请退款/售后
    alt 正常流程
        O->>L: 11. 创建物流单
        L->>DB: 插入 Logistics (Status=0)
        L-->>O: 物流单创建成功
        O->>DB: Update Order: Status=3
        O-->>C: 发货通知
        C-->>U: 订单已发货
        
        U->>C: 12. 确认收货
        C->>O: 13. 确认收货
        O->>DB: Update Order: Status=4
        O->>S: 14. 扣减库存
        S->>DB: Update Stock:<br/>LockedStock -= qty<br/>TotalStock -= qty
        S-->>O: 扣减成功
        O->>DB: Update Logistics:<br/>Status=3, ReceiveAt
        O-->>C: 订单完成
        C-->>U: 交易完成
    else 取消/退款流程
        O->>S: 回滚库存
        S->>DB: Update Stock:<br/>恢复AvailableStock<br/>清除Frozen/Locked
        S-->>O: 回滚成功
        O->>DB: Update Order: Status=5/6
        O-->>C: 订单状态更新
        C-->>U: 取消/退款成功
    end
    end
```

### 1.3 订单数据流图

```mermaid
flowchart TD
    Start([用户下单]) --> CheckStock{检查库存<br/>Stock.AvailableStock}
    
    CheckStock -->|库存不足| Error1[返回库存不足]
    CheckStock -->|库存充足| CreateOrder[创建订单<br/>Order.Status=1]
    
    CreateOrder --> FreezeStock[冻结库存<br/>Stock.FrozenStock+=qty]
    FreezeStock --> WaitPay[等待支付<br/>订单有效期15分钟]
    
    WaitPay -->|超时未支付| CancelOrder[自动取消订单]
    CancelOrder --> UnfreezeStock[解冻库存<br/>Stock.AvailableStock+=qty]
    UnfreezeStock --> End1([订单取消])
    
    WaitPay -->|用户支付| LockStock[锁定库存<br/>Stock.LockedStock+=qty]
    LockStock --> UpdateOrder[更新订单<br/>Order.Status=2]
    UpdateOrder --> CreateLogistics[创建物流单]
    
    CreateLogistics --> Ship[发货<br/>Order.Status=3]
    Ship --> Receive[用户收货<br/>Order.Status=4]
    Receive --> DeductStock[扣减库存<br/>Stock.TotalStock-=qty]
    DeductStock --> Complete([订单完成])
    
    WaitPay -->|用户主动取消| CancelOrder2[取消订单]
    CancelOrder2 --> UnfreezeStock2[解冻库存]
    UnfreezeStock2 --> End2([订单取消])
```

---

## 二、库存全流程

### 2.1 库存状态定义

基于 `Stock` 模型的库存状态：

| 字段 | 说明 | 计算关系 |
|------|------|----------|
| TotalStock | 总库存 | 物理库存总量 |
| AvailableStock | 可用库存 | 可售卖的库存 |
| FrozenStock | 冻结库存 | 已下单未支付，占用中 |
| LockedStock | 锁定库存 | 已支付，待发货 |

**库存公式**: `TotalStock = AvailableStock + FrozenStock + LockedStock`

### 2.2 库存时序图

```mermaid
sequenceDiagram
    actor U as 用户
    participant API as API Gateway
    participant SS as 秒杀服务
    participant IS as 库存服务
    participant RS as Redis
    participant DB as MySQL

    %% 1. 库存查询
    rect rgb(240, 248, 255)
    Note over U,DB: 场景1: 库存查询
    U->>API: 1. 查询商品库存
    API->>IS: 2. 获取库存
    IS->>RS: 3. 读取缓存库存
    alt 缓存命中
        RS-->>IS: 返回缓存数据
    else 缓存未命中
        IS->>DB: 4. 查询 Stock 表
        DB-->>IS: 返回 Stock 记录
        IS->>RS: 5. 写入缓存
    end
    IS-->>API: 返回 AvailableStock
    API-->>U: 显示可售数量
    end

    %% 2. 秒杀下单 - 库存预占
    rect rgb(255, 248, 240)
    Note over U,DB: 场景2: 秒杀库存预占
    U->>API: 6. 秒杀下单请求
    API->>SS: 7. 验证秒杀资格
    SS->>IS: 8. 预占库存
    
    IS->>RS: 9. Lua脚本扣减<br/>DECR AvailableStock<br/>INCR FrozenStock
    RS-->>IS: 返回扣减结果
    
    alt 扣减成功
        IS->>DB: 10. 异步同步数据库
        IS-->>SS: 预占成功
        SS-->>API: 返回下单Token
        API-->>U: 进入下单流程
    else 库存不足
        IS-->>SS: 库存不足
        SS-->>API: 返回秒杀结束
        API-->>U: 显示已售罄
    end
    end

    %% 3. 支付成功 - 库存锁定
    rect rgb(240, 255, 240)
    Note over U,DB: 场景3: 支付后库存锁定
    U->>API: 11. 完成支付
    API->>SS: 12. 支付回调
    SS->>IS: 13. 确认库存
    
    IS->>RS: 14. Lua脚本转移<br/>DECR FrozenStock<br/>INCR LockedStock
    RS-->>IS: 转移成功
    
    IS->>DB: 15. 同步更新 Stock
    IS-->>SS: 确认成功
    SS-->>API: 订单状态更新
    API-->>U: 支付完成通知
    end

    %% 4. 发货完成 - 库存扣减
    rect rgb(255, 240, 255)
    Note over U,DB: 场景4: 发货后库存扣减
    U->>API: 16. 确认收货
    API->>IS: 17. 完成库存扣减
    
    IS->>RS: 18. 扣减库存<br/>DECR LockedStock<br/>DECR TotalStock
    IS->>DB: 19. 更新 Stock<br/>扣减总库存
    IS-->>API: 扣减完成
    API-->>U: 交易完成
    end

    %% 5. 取消订单 - 库存回滚
    rect rgb(255, 255, 240)
    Note over U,DB: 场景5: 订单取消库存回滚
    U->>API: 20. 取消订单/超时
    API->>IS: 21. 回滚库存
    
    alt 未支付取消
        IS->>RS: 22. 解冻库存<br/>DECR FrozenStock<br/>INCR AvailableStock
    else 已支付取消
        IS->>RS: 23. 解锁库存<br/>DECR LockedStock<br/>INCR AvailableStock
    end
    
    IS->>DB: 24. 同步回滚
    IS-->>API: 回滚完成
    API-->>U: 取消成功
    end
```

### 2.3 库存数据流图

```mermaid
flowchart TD
    subgraph 库存状态机
    A[总库存<br/>TotalStock] 
    B[可用库存<br/>AvailableStock]
    C[冻结库存<br/>FrozenStock]
    D[锁定库存<br/>LockedStock]
    end

    Start([用户访问商品]) --> QueryStock{查询库存}
    QueryStock -->|缓存查询| Redis[(Redis<br/>库存缓存)]
    QueryStock -->|DB查询| DB[(MySQL<br/>Stock表)]
    
    Redis --> Display[展示可售数量]
    DB --> Display
    
    Display --> SecKill{秒杀下单}
    SecKill -->|库存充足| Freeze[冻结库存<br/>Available→Frozen]
    SecKill -->|库存不足| SoldOut([已售罄])
    
    Freeze --> WaitPay[等待支付]
    WaitPay -->|支付成功| Lock[锁定库存<br/>Frozen→Locked]
    WaitPay -->|超时/取消| Unfreeze[解冻库存<br/>Frozen→Available]
    
    Lock --> WaitShip[等待发货]
    WaitShip -->|用户收货| Deduct[扣减库存<br/>Locked--<br/>Total--]
    WaitShip -->|取消订单| Unlock[解锁库存<br/>Locked→Available]
    
    Deduct --> Complete([交易完成])
    Unfreeze --> Cancelled([订单取消])
    Unlock --> Cancelled2([订单取消])
    SoldOut --> End([结束])
```

---

## 三、积分全流程

### 3.1 积分类型定义

基于 `Point` 模型的 `Type` 字段：

| 类型码 | 类型名称 | 触发场景 |
|--------|----------|----------|
| 1 | 获取 | 购物返积分、签到、活动奖励 |
| 2 | 消费 | 积分抵扣订单金额 |
| 3 | 过期 | 积分到期自动清零 |
| 4 | 退款 | 订单退款返还积分 |

### 3.2 积分时序图

```mermaid
sequenceDiagram
    actor U as 用户
    participant API as API Gateway
    participant OS as 订单服务
    participant PS as 积分服务
    participant DB as MySQL

    %% 1. 积分获取 - 购物返积分
    rect rgb(240, 248, 255)
    Note over U,DB: 场景1: 购物返积分
    U->>API: 1. 确认收货
    API->>OS: 2. 订单完成
    OS->>PS: 3. 计算返积分<br/>(OrderItem.Point)
    
    PS->>DB: 4. 插入 Point 流水<br/>Type=1, Source=购物
    PS->>DB: 5. 更新 MemberPoint<br/>TotalPoint += point<br/>AvailablePoint += point
    
    PS-->>OS: 返积分成功
    OS-->>API: 订单完成
    API-->>U: 获得X积分通知
    end

    %% 2. 积分消费
    rect rgb(255, 248, 240)
    Note over U,DB: 场景2: 积分消费抵扣
    U->>API: 6. 创建订单(使用积分)
    API->>PS: 7. 检查积分余额<br/>MemberPoint.AvailablePoint
    PS->>DB: 查询积分余额
    DB-->>PS: 返回 AvailablePoint
    
    alt 积分充足
        PS-->>API: 余额充足
        API->>OS: 8. 创建订单
        OS->>PS: 9. 冻结积分
        
        PS->>DB: 10. 插入 Point 流水<br/>Type=2(预扣)
        PS->>DB: 11. 更新 MemberPoint<br/>AvailablePoint -= point
        PS-->>OS: 冻结成功
        OS-->>API: 订单创建成功<br/>PointAmount=X
        API-->>U: 订单创建成功
        
        U->>API: 12. 完成支付
        API->>PS: 13. 确认积分消费
        PS->>DB: 14. 更新 Point<br/>Type=2(确认消费)
        PS-->>API: 积分扣减完成
    else 积分不足
        PS-->>API: 余额不足
        API-->>U: 提示积分不足
    end
    end

    %% 3. 积分过期
    rect rgb(240, 255, 240)
    Note over U,DB: 场景3: 积分过期处理
    Note right of PS: 定时任务扫描<br/>Point.ExpireAt < now
    PS->>DB: 15. 查询过期积分
    DB-->>PS: 返回过期记录
    
    PS->>DB: 16. 插入 Point 流水<br/>Type=3, Point=-expired
    PS->>DB: 17. 更新 MemberPoint<br/>TotalPoint -= expired<br/>ExpiredPoint += expired<br/>AvailablePoint -= expired
    PS-->>PS: 发送过期通知
    end

    %% 4. 积分退款
    rect rgb(255, 240, 255)
    Note over U,DB: 场景4: 订单退款返积分
    U->>API: 18. 申请退款
    API->>OS: 19. 处理退款
    OS->>PS: 20. 返还积分
    
    PS->>DB: 21. 插入 Point 流水<br/>Type=4, Source=退款
    PS->>DB: 22. 更新 MemberPoint<br/>TotalPoint += refund<br/>UsedPoint -= refund<br/>AvailablePoint += refund
    
    PS-->>OS: 返积分成功
    OS-->>API: 退款完成
    API-->>U: 退款成功通知
    end
```

### 3.3 积分数据流图

```mermaid
flowchart TD
    subgraph 积分账户
    MP[MemberPoint<br/>会员积分总表]
    MP_Total[TotalPoint<br/>总积分]
    MP_Used[UsedPoint<br/>已用积分]
    MP_Expired[ExpiredPoint<br/>已过期]
    MP_Avail[AvailablePoint<br/>可用积分]
    end

    subgraph 积分流水
    P[Point<br/>积分流水表]
    P_Type[Type: 1获取/2消费/3过期/4退款]
    P_Balance[Balance: 变动后余额]
    end

    Start([积分变动]) --> CheckType{积分类型}
    
    CheckType -->|Type=1 获取| Earn[获得积分]
    Earn --> UpdateMP1[更新 MemberPoint:<br/>Total++, Available++]
    UpdateMP1 --> InsertP1[插入 Point:<br/>Type=1, Balance=new]
    
    CheckType -->|Type=2 消费| Spend[消费积分]
    Spend --> CheckBalance{检查可用积分}
    CheckBalance -->|充足| UpdateMP2[更新 MemberPoint:<br/>Used++, Available--]
    UpdateMP2 --> InsertP2[插入 Point:<br/>Type=2, Balance=new]
    CheckBalance -->|不足| Error[返回余额不足]
    
    CheckType -->|Type=3 过期| Expire[积分过期]
    Expire --> UpdateMP3[更新 MemberPoint:<br/>Expired++, Available--]
    UpdateMP3 --> InsertP3[插入 Point:<br/>Type=3, Balance=new]
    
    CheckType -->|Type=4 退款| Refund[积分退款]
    Refund --> UpdateMP4[更新 MemberPoint:<br/>Used--, Available++]
    UpdateMP4 --> InsertP4[插入 Point:<br/>Type=4, Balance=new]
    
    InsertP1 --> Notify[发送积分变动通知]
    InsertP2 --> Notify
    InsertP3 --> Notify
    InsertP4 --> Notify
    
    Notify --> End([完成])
    Error --> End2([失败])
```

---

## 四、数据流程图

### 4.1 系统整体数据流

```mermaid
flowchart TD
    subgraph 客户端层
    Web[Web/App 客户端]
    end

    subgraph 网关层
    Gateway[API Gateway]
    Router[路由分发]
    end

    subgraph 服务层
    OS[订单服务<br/>Order]
    SS[秒杀服务<br/>Seckill]
    IS[库存服务<br/>Stock]
    PS[积分服务<br/>Point]
    LS[物流服务<br/>Logistics]
    US[用户服务<br/>Member]
    end

    subgraph 数据层
    MySQL[(MySQL<br/>业务数据)]
    Redis[(Redis<br/>缓存/计数)]
    end

    subgraph 外部服务
    Payment[支付服务]
    Message[消息队列]
    end

    %% 数据流向
    Web -->|HTTP/gRPC| Gateway
    Gateway --> Router
    
    Router -->|订单相关| OS
    Router -->|秒杀相关| SS
    Router -->|库存查询| IS
    Router -->|积分查询| PS
    Router -->|用户信息| US
    
    OS -->|读写| MySQL
    OS -->|缓存| Redis
    OS -->|扣减库存| IS
    OS -->|扣减积分| PS
    OS -->|创建物流| LS
    OS -->|发起支付| Payment
    
    SS -->|库存预占| IS
    SS -->|生成订单| OS
    
    IS -->|读写| MySQL
    IS -->|缓存| Redis
    
    PS -->|读写| MySQL
    
    LS -->|读写| MySQL
    
    US -->|读写| MySQL
    
    OS -->|订单状态变更| Message
    PS -->|积分变动| Message
    IS -->|库存预警| Message
    
    Message -->|通知| Web
```

### 4.2 核心业务流程数据流

```mermaid
flowchart LR
    subgraph 秒杀流程
    A1[用户参与秒杀] --> B1[库存预占<br/>Redis Lua]
    B1 --> C1[创建订单<br/>Order.Status=1]
    C1 --> D1[冻结库存<br/>Stock.FrozenStock++]
    D1 --> E1[等待支付]
    E1 -->|支付成功| F1[锁定库存<br/>Stock.LockedStock++]
    F1 --> G1[发货扣减<br/>Stock.TotalStock--]
    E1 -->|超时/取消| H1[回滚库存<br/>AvailableStock++]
    end

    subgraph 积分流程
    A2[订单完成] --> B2[计算返积分<br/>OrderItem.Point]
    B2 --> C2[更新 MemberPoint<br/>Available++]
    C2 --> D2[插入 Point<br/>Type=1]
    
    E2[使用积分] --> F2[检查余额<br/>AvailablePoint]
    F2 --> G2[冻结积分<br/>Point.Type=2]
    G2 --> H2[支付确认]
    H2 --> I2[扣减积分<br/>Used++]
    end
```

---

## 五、核心实体关系

### 5.1 数据库实体关系图

```mermaid
erDiagram
    MEMBER ||--o{ ORDER : places
    MEMBER ||--o{ MEMBER_POINT : has
    MEMBER ||--o{ POINT : earns
    
    ORDER ||--|{ ORDER_ITEM : contains
    ORDER ||--o| LOGISTICS : shipped_by
    ORDER ||--o{ POINT : generates
    
    PRODUCT ||--o{ STOCK : has
    PRODUCT ||--o{ SKU : has_variants
    PRODUCT ||--o{ ORDER_ITEM : ordered_in
    
    STOCK }|--|| SKU : tracks
    
    MEMBER_POINT {
        uint MemberID PK
        int TotalPoint
        int UsedPoint
        int ExpiredPoint
        int AvailablePoint
    }
    
    POINT {
        uint ID PK
        uint MemberID FK
        uint OrderID FK
        string PointNo
        int Type
        int Point
        int Balance
        string Source
        time ExpireAt
    }
    
    ORDER {
        uint ID PK
        string OrderNo
        uint MemberID FK
        int Status
        float TotalAmount
        float PayAmount
        float PointAmount
        int PaymentMethod
    }
    
    ORDER_ITEM {
        uint ID PK
        uint OrderID FK
        uint ProductID FK
        uint SkuID FK
        int Quantity
        float Price
        int Point
    }
    
    STOCK {
        uint ID PK
        uint ProductID FK
        uint SkuID FK
        int TotalStock
        int AvailableStock
        int FrozenStock
        int LockedStock
    }
    
    PRODUCT {
        uint ID PK
        string Name
        float Price
        int Status
    }
    
    SKU {
        uint ID PK
        uint ProductID FK
        string SkuNo
        string Specs
        float Price
        int Stock
    }
    
    LOGISTICS {
        uint ID PK
        uint OrderID FK
        string LogisticsNo
        int Status
        time ShipAt
        time ReceiveAt
    }
    
    MEMBER {
        uint ID PK
        string Username
        string Mobile
        int Status
    }
```

### 5.2 关键字段说明

#### Order 订单表
- `OrderNo`: 唯一订单号
- `Status`: 1待付款/2待发货/3待收货/4已完成/5已取消/6售后
- `TotalAmount`: 订单总额
- `PayAmount`: 实付金额
- `PointAmount`: 积分抵扣金额
- `PaymentMethod`: 1微信/2支付宝/3积分

#### Stock 库存表
- `TotalStock`: 总库存数量
- `AvailableStock`: 可用库存(可售卖)
- `FrozenStock`: 冻结库存(已下单未支付)
- `LockedStock`: 锁定库存(已支付待发货)

#### Point 积分流水表
- `Type`: 1获取/2消费/3过期/4退款
- `Point`: 变动积分数(正数增加,负数减少)
- `Balance`: 变动后的积分余额
- `ExpireAt`: 积分过期时间

#### MemberPoint 会员积分总表
- `TotalPoint`: 累计获得积分
- `UsedPoint`: 已使用积分
- `ExpiredPoint`: 已过期积分
- `AvailablePoint`: 当前可用积分

---

## 附录：模型定义源码

### Order 订单模型
```go
type Order struct {
    ID             uint       `gorm:"primaryKey;comment:'订单ID'"`
    OrderNo        string     `gorm:"type:varchar(50);uniqueIndex;not null;comment:'订单号'"`
    MemberID       uint       `gorm:"index;not null;comment:'会员ID'"`
    Status         int        `gorm:"type:tinyint(2);default:1;comment:'状态:1待付款2待发货3待收货4已完成5已取消6售后'"`
    TotalAmount    float64    `gorm:"type:decimal(10,2);not null;comment:'订单总额'"`
    PayAmount      float64    `gorm:"type:decimal(10,2);not null;comment:'实付金额'"`
    PointAmount    float64    `gorm:"type:decimal(10,2);default:0;comment:'积分抵扣'"`
    PaymentMethod  int        `gorm:"type:tinyint(1);comment:'支付方式:1微信2支付宝3积分'"`
    CreatedAt      time.Time  `gorm:"comment:'下单时间'"`
}
```

### Stock 库存模型
```go
type Stock struct {
    ID             uint       `gorm:"primaryKey;comment:'库存ID'"`
    ProductID      uint       `gorm:"uniqueIndex;not null;comment:'商品ID'"`
    TotalStock     int        `gorm:"default:0;comment:'总库存'"`
    AvailableStock int        `gorm:"default:0;comment:'可用库存'"`
    FrozenStock    int        `gorm:"default:0;comment:'冻结库存(已下单未支付)'"`
    LockedStock    int        `gorm:"default:0;comment:'锁定库存(已支付)'"`
}
```

### Point 积分模型
```go
type Point struct {
    ID        uint       `gorm:"primaryKey;comment:'积分ID'"`
    MemberID  uint       `gorm:"index;not null;comment:'会员ID'"`
    OrderID   uint       `gorm:"index;comment:'订单ID'"`
    Type      int        `gorm:"type:tinyint(1);comment:'类型:1获取2消费3过期4退款'"`
    Point     int        `gorm:"not null;comment:'积分数值'"`
    Balance   int        `gorm:"default:0;comment:'积分余额'"`
    ExpireAt  *time.Time `gorm:"type:datetime;comment:'过期时间'"`
}
```

---

**文档结束**
**记录时间**: 基于当前代码库模型定义分析
**数据来源**: server/models/model.go
