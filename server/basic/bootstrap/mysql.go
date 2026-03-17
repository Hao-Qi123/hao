package bootstrap

import (
	"fmt"
	"seckil/server/basic/config"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var once sync.Once

func InitMysql() {
	var err error
	configMysql := config.GlobalConfig.Mysql
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		configMysql.User,
		configMysql.Password,
		configMysql.Host,
		configMysql.Port,
		configMysql.Database)
	once.Do(func() {
		config.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("数据库连接失败")
		}
		fmt.Println("数据库连接成功")

		sqlDB, err := config.DB.DB()
		if err != nil {
			panic("数据池连接失败")
		}

		// SetMaxIdleConns 设置空闲连接池中连接的最大数量。
		sqlDB.SetMaxIdleConns(10)

		// SetMaxOpenConns 设置打开数据库连接的最大数量。
		sqlDB.SetMaxOpenConns(100)

		// SetConnMaxLifetime 设置了可以重新使用连接的最大时间。
		sqlDB.SetConnMaxLifetime(time.Hour)
		fmt.Println("数据池连接成功")

	})
	//err = config.DB.AutoMigrate(
	//	&models.Member{},
	//	&models.Product{},
	//	&models.Stock{},
	//	&models.Point{},
	//	&models.Order{},
	//	&models.OrderItem{},
	//	&models.Logistics{})
	//if err != nil {
	//	panic("数据库迁移失败")
	//}
	//fmt.Println("数据库迁移成功")
}
