package handler

import (
	sqlc "crud/db/sqlc"
	sqlx "crud/db/sqlx"
	"crud/pkg/response"
	"database/sql"
	"fmt"

	"github.com/labstack/echo/v4"
)

// author 全局作者处理器实例
var author *AuthorApi

// AuthorApi 作者API处理结构体
type AuthorApi struct {
	*BaseCrudHandler[sqlc.Author, sqlc.AuthorUpdate]
}

// NewAuthorApi 创建新的作者API处理器
func NewAuthorApi(dbConn *sql.DB) *AuthorApi {
	crud := sqlx.NewModelWithGlobal[sqlc.Author](dbConn)
	return &AuthorApi{
		BaseCrudHandler: NewBaseCrudHandler[sqlc.Author, sqlc.AuthorUpdate]("作者", crud),
	}
}

type AuthorWithBooks struct {
	*sqlc.Author
	Books []sqlc.Book `json:"books"`
}

func (h *AuthorApi) GetAuthorWithBooks(c echo.Context) error {
	var singleId SingleId
	if err := c.Bind(&singleId); err != nil {
		return response.BadRequest(err)
	}
	if singleId.Id == 0 {
		return response.BadRequest(fmt.Errorf("作者ID不能为空"))
	}
	authors, err := h.BaseCrudHandler.crud.Sqlc.GetAuthorWithBooks(c.Request().Context(), singleId.Id)
	if err != nil {
		return response.DatabaseError(err)
	}
	if len(authors) == 0 {
		return response.NotFound(fmt.Errorf("作者不存在"))
	}
	var result AuthorWithBooks
	result.Author = &sqlc.Author{
		ID:   authors[0].AuthorID,
		Name: authors[0].AuthorName,
		Bio:  authors[0].AuthorBio,
	}
	var books []sqlc.Book
	for _, v := range authors {
		books = append(books, sqlc.Book{
			ID:    v.BookID.Int64,
			Title: v.BookTitle.String,
		})
	}
	result.Books = books
	return response.Success(c, result)
}
