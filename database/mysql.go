package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB = InitDB()

// 注意这里：返回值改为 (*gorm.DB, error)
func InitDB() *gorm.DB {
	dsn := "blog_user:blog_user@tcp(127.0.0.1:3306)/go_blog?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接 MySQL 失败: %v", err)
	}

	// 检查数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取数据库连接失败: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("MySQL Ping 失败: %v", err)
	}

	log.Println("MySQL 连接成功")
	return db
}

func CloseMySQL() error {

	if DB == nil {
		return nil
	}

	// 从 GORM 获取底层 *sql.DB
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	// 关闭数据库连接池
	return sqlDB.Close()
}
