package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// NewMySQLConnector 创建一个新的 MySQL 连接器
func NewMySQLConnector(config DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}

	// 设置连接池参数
	db.SetMaxIdleConns(10)           // 最大空闲连接数
	db.SetMaxOpenConns(100)          // 最大打开连接数
	db.SetConnMaxLifetime(time.Hour) // 连接最大生命周期

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping 数据库失败: %v", err)
	}

	return db, nil
}
