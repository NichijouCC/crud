package handler

import (
	"crud/db/sqlx"
	"crud/pkg/response"
	"fmt"

	"github.com/labstack/echo/v4"
)

// BaseHandler 基础CRUD处理器接口
type BaseHandler[T any, U any] interface {
	GetById(c echo.Context) error
	GetByIds(c echo.Context) error
	GetByFilter(c echo.Context) error
	GetAll(c echo.Context) error
	Create(c echo.Context) error
	DeleteById(c echo.Context) error
	DeleteByIds(c echo.Context) error
	DeleteByFilter(c echo.Context) error
	UpdateById(c echo.Context) error
	UpdateByIds(c echo.Context) error
	UpdateByFilter(c echo.Context) error
}

// BaseCrudHandler 基础CRUD处理器实现
type BaseCrudHandler[T sqlx.ITable, U sqlx.ITableUpdate] struct {
	resourceName string
	crud         *sqlx.Table[T]
}

// NewBaseCrudHandler 创建新的基础CRUD处理器
func NewBaseCrudHandler[T sqlx.ITable, U sqlx.ITableUpdate](resourceName string, crud *sqlx.Table[T]) *BaseCrudHandler[T, U] {
	return &BaseCrudHandler[T, U]{
		resourceName: resourceName,
		crud:         crud,
	}
}

// GetById 根据ID获取单个资源
func (h *BaseCrudHandler[T, U]) GetById(c echo.Context) error {
	var singleId SingleId
	if err := c.Bind(&singleId); err != nil {
		return response.BadRequest(err)
	}
	if singleId.Id == 0 {
		return response.BadRequest(fmt.Errorf("%s ID不能为空", h.resourceName))
	}
	item, err := h.crud.FindOneById(c.Request().Context(), singleId.Id)
	if err != nil {
		return response.DatabaseError(err)
	}
	return response.Success(c, item)
}

// GetByIds 根据ID列表获取多个资源
func (h *BaseCrudHandler[T, U]) GetByIds(c echo.Context) error {
	var groupIds GroupIds
	if err := c.Bind(&groupIds); err != nil {
		return response.BadRequest(err)
	}
	if len(groupIds.Ids) == 0 {
		return response.BadRequest(fmt.Errorf("%s ID列表不能为空", h.resourceName))
	}
	items, err := h.crud.FindSomeByIds(c.Request().Context(), groupIds.Ids)
	if err != nil {
		return response.DatabaseError(err)
	}
	return response.Success(c, items)
}

// GetByFilter 根据过滤条件获取资源
func (h *BaseCrudHandler[T, U]) GetByFilter(c echo.Context) error {
	params := c.QueryParams()
	filter, err := sqlx.ParseQueryFilterFromUrlParams(params)
	if err != nil {
		return response.BadRequest(err)
	}
	items, err := h.crud.FindSomeByFilter(c.Request().Context(), filter)
	if err != nil {
		return response.DatabaseError(err)
	}
	return response.Success(c, items)
}

// GetAll 获取所有资源
func (h *BaseCrudHandler[T, U]) GetAll(c echo.Context) error {
	return h.GetByFilter(c)
}

// Create 创建新资源
func (h *BaseCrudHandler[T, U]) Create(c echo.Context) error {
	var item T
	if err := c.Bind(&item); err != nil {
		return response.BadRequest(err)
	}
	err := h.crud.CreateOne(c.Request().Context(), item)
	if err != nil {
		return response.DatabaseError(err)
	}
	return response.Success(c, item)
}

// DeleteById 删除单个资源
func (h *BaseCrudHandler[T, U]) DeleteById(c echo.Context) error {
	var singleId SingleId
	if err := c.Bind(&singleId); err != nil {
		return response.BadRequest(err)
	}
	if singleId.Id == 0 {
		return response.BadRequest(fmt.Errorf("%s ID不能为空", h.resourceName))
	}
	err := h.crud.DeleteOne(c.Request().Context(), singleId.Id)
	if err != nil {
		return response.DatabaseError(err)
	}
	return response.Success(c, nil)
}

// DeleteByIds 批量删除资源
func (h *BaseCrudHandler[T, U]) DeleteByIds(c echo.Context) error {
	var groupIds GroupIds
	if err := c.Bind(&groupIds); err != nil {
		return response.BadRequest(err)
	}
	if len(groupIds.Ids) == 0 {
		return response.BadRequest(fmt.Errorf("%s ID列表不能为空", h.resourceName))
	}
	err := h.crud.DeleteSomeByIds(c.Request().Context(), groupIds.Ids)
	if err != nil {
		return response.DatabaseError(err)
	}
	return response.Success(c, nil)
}

// DeleteByFilter 根据过滤条件删除资源
func (h *BaseCrudHandler[T, U]) DeleteByFilter(c echo.Context) error {
	params := c.QueryParams()
	filter, err := sqlx.ParseQueryFilterFromUrlParams(params)
	if err != nil {
		return response.BadRequest(err)
	}
	err = h.crud.DeleteSomeByFilter(c.Request().Context(), filter)
	if err != nil {
		return response.DatabaseError(err)
	}
	return response.Success(c, nil)
}

// UpdateById 更新单个资源
func (h *BaseCrudHandler[T, U]) UpdateById(c echo.Context) error {
	var item U
	if err := c.Bind(&item); err != nil {
		return response.BadRequest(err)
	}
	if item.GetId() == 0 {
		return response.BadRequest(fmt.Errorf("%s ID不能为空", h.resourceName))
	}
	err := h.crud.UpdateOne(c.Request().Context(), item)
	if err != nil {
		return response.DatabaseError(err)
	}
	return response.Success(c, item)
}

// UpdateByIds 批量更新资源
func (h *BaseCrudHandler[T, U]) UpdateByIds(c echo.Context) error {
	var groupIds GroupIds
	if err := c.Bind(&groupIds); err != nil {
		return response.BadRequest(err)
	}
	if len(groupIds.Ids) == 0 {
		return response.BadRequest(fmt.Errorf("%s ID列表不能为空", h.resourceName))
	}
	var item U
	if err := c.Bind(&item); err != nil {
		return response.BadRequest(err)
	}
	err := h.crud.UpdateSomeByIds(c.Request().Context(), item, groupIds.Ids)
	if err != nil {
		return response.DatabaseError(err)
	}
	return response.Success(c, item)
}

// UpdateByFilter 根据过滤条件更新资源
func (h *BaseCrudHandler[T, U]) UpdateByFilter(c echo.Context) error {
	params := c.QueryParams()
	filter, err := sqlx.ParseQueryFilterFromUrlParams(params)
	if err != nil {
		return response.BadRequest(err)
	}
	var item U
	if err := c.Bind(&item); err != nil {
		return response.BadRequest(err)
	}
	err = h.crud.UpdateSomeByFilter(c.Request().Context(), item, filter)
	if err != nil {
		return response.DatabaseError(err)
	}
	return response.Success(c, item)
}

type GroupIds struct {
	Ids []int64 `json:"ids" param:"ids" query:"ids" form:"ids"`
}

type SingleId struct {
	Id int64 `json:"id" param:"id" query:"id" form:"id"`
}
