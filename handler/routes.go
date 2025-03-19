package handler

import (
	"database/sql"

	"github.com/labstack/echo/v4"
)

// RegisterRoutes 注册所有API路由
func RegisterRoutes(e *echo.Echo, dbConn *sql.DB) {
	// 初始化处理器
	author = NewAuthorApi(dbConn)
	book = NewBookApi(dbConn)

	e.GET("/ping", PingHandler)

	// 作者相关路由
	e.GET("/authors", author.GetAll)                       // 获取所有作者
	e.GET("/authors/:ids", author.GetByIds)                // 获取多个作者
	e.POST("/authors", author.Create)                      // 创建作者
	e.DELETE("/authors/:id", author.DeleteById)            // 删除单个作者
	e.PUT("/authors/:id", author.UpdateById)               // 更新作者信息
	e.GET("/authors/:id/books", author.GetAuthorWithBooks) // 获取作者及其书籍

	// 书籍相关路由
	e.GET("/books", book.GetAll)            // 获取所有书籍
	e.GET("/books/:ids", book.GetByIds)     // 获取多个书籍
	e.POST("/books", book.Create)           // 创建书籍
	e.DELETE("/books/:id", book.DeleteById) // 删除单个书籍
	e.PUT("/books/:id", book.UpdateById)    // 更新书籍信息
}
