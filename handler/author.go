package handler

import (
	"crud/db/sqlx"
	"crud/pkg/response"
	"database/sql"
	"fmt"

	"github.com/labstack/echo/v4"
)

// AuthorApi 作者API处理结构体
type AuthorApi struct {
	crud *sqlx.Model[sqlx.Author] // 作者CRUD操作模型
}

// NewAuthorApi 创建新的作者API处理器
func NewAuthorApi(db *sql.DB) *AuthorApi {
	crud := sqlx.NewModel[sqlx.Author](db)
	return &AuthorApi{crud: crud}
}

// RegisterAuthorRoutes 初始化作者相关的API路由
func RegisterAuthorRoutes(e *echo.Echo, db *sql.DB) {
	api := NewAuthorApi(db)
	e.GET("/authors", api.GetAuthors)          // 获取所有作者
	e.GET("/authors/:ids", api.GetAuthorByIds) // 获取多个作者
	e.POST("/authors", api.CreateAuthor)       // 创建作者
	e.DELETE("/authors/:id", api.DeleteAuthor) // 删除单个作者
	e.PUT("/authors/:id", api.UpdateAuthor)    // 更新作者信息
}

// GetAuthors 获取所有作者信息
func (a *AuthorApi) GetAuthors(c echo.Context) error {
	params := c.QueryParams()
	filter, err := sqlx.ParseFilterFromContext(params)
	if err != nil {
		return response.BadRequest(err)
	}
	var authors []*sqlx.Author
	if filter == nil {
		authors, err = a.crud.FindRows(c.Request().Context())
	} else {
		authors, err = a.crud.FindRowsByFilter(c.Request().Context(), filter)
	}
	if err != nil {
		return response.DatabaseError(err)
	}
	return response.Success(c, authors)
}

func (a *AuthorApi) GetAuthorByIds(c echo.Context) error {
	var groupIds GroupIds
	if err := c.Bind(&groupIds); err != nil {
		return response.BadRequest(err)
	}
	if len(groupIds.Ids) == 0 {
		return response.BadRequest(fmt.Errorf("作者ID列表不能为空"))
	}
	authors, err := a.crud.FindRowsByIDs(c.Request().Context(), groupIds.Ids)
	if err != nil {
		return response.DatabaseError(err)
	}
	return response.Success(c, authors)
}

func (a *AuthorApi) GetGroupWithFilter(c echo.Context) error {
	var filter sqlx.Filter
	if err := c.Bind(&filter); err != nil {
		return response.BadRequest(err)
	}
	authors, err := a.crud.FindRowsByFilter(c.Request().Context(), &filter)
	if err != nil {
		return response.DatabaseError(err)
	}
	return response.Success(c, authors)
}

// CreateAuthor 创建新作者
func (a *AuthorApi) CreateAuthor(c echo.Context) error {
	var author sqlx.Author
	if err := c.Bind(&author); err != nil {
		return response.BadRequest(err)
	}
	err := a.crud.Create(c.Request().Context(), &author)
	if err != nil {
		return response.DatabaseError(err)
	}
	return response.Success(c, author)
}

// DeleteAuthor 删除单个作者
func (a *AuthorApi) DeleteAuthor(c echo.Context) error {
	var singleId SingleId
	if err := c.Bind(&singleId); err != nil {
		return response.BadRequest(err)
	}
	if singleId.Id == 0 {
		return response.BadRequest(fmt.Errorf("作者ID不能为空"))
	}
	err := a.crud.Delete(c.Request().Context(), singleId.Id)
	if err != nil {
		return response.DatabaseError(err)
	}
	return response.Success(c, nil)
}

// DeleteAuthors 批量删除作者
func (a *AuthorApi) DeleteAuthors(c echo.Context) error {
	var groupIds GroupIds
	if err := c.Bind(&groupIds); err != nil {
		return response.BadRequest(err)
	}
	if len(groupIds.Ids) == 0 {
		return response.BadRequest(fmt.Errorf("作者ID列表不能为空"))
	}
	err := a.crud.DeleteRowsByIDs(c.Request().Context(), groupIds.Ids)
	if err != nil {
		return response.DatabaseError(err)
	}
	return response.Success(c, nil)
}

// UpdateAuthor 更新作者信息
func (a *AuthorApi) UpdateAuthor(c echo.Context) error {
	author := new(sqlx.AuthorUpdate)
	if err := c.Bind(author); err != nil {
		return response.BadRequest(err)
	}
	if author.Id == 0 {
		return response.BadRequest(fmt.Errorf("作者ID不能为空"))
	}
	err := a.crud.Update(c.Request().Context(), author)
	if err != nil {
		return response.DatabaseError(err)
	}
	return response.Success(c, author)
}
