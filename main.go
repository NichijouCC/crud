package main

import (
	"crud/handler"
	"crud/middleware"
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func main() {
	// 连接数据库
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/crud?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 创建 Echo 实例
	e := echo.New()

	// 添加中间件
	e.Use(middleware.Recover())      // 恢复中间件
	e.Use(middleware.CORS())         // 跨域中间件
	e.Use(middleware.Logger())       // 日志中间件
	e.Use(middleware.ErrorHandler()) // 错误处理中间件

	// 注册路由
	handler.RegisterAuthorRoutes(e, db)

	// 启动服务器
	e.Logger.Fatal(e.Start(":8080"))
}
