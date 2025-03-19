package handler

import (
	sqlc "crud/db/sqlc"
	sqlx "crud/db/sqlx"
	"database/sql"
)

// book 全局书籍处理器实例
var book *BookApi

// BookApi 书籍API处理结构体
type BookApi struct {
	*BaseCrudHandler[sqlc.Book, sqlc.BookUpdate]
}

// NewBookApi 创建新的书籍API处理器
func NewBookApi(dbConn *sql.DB) *BookApi {
	crud := sqlx.NewModelWithGlobal[sqlc.Book](dbConn)
	return &BookApi{
		BaseCrudHandler: NewBaseCrudHandler[sqlc.Book, sqlc.BookUpdate]("书籍", crud),
	}
}
