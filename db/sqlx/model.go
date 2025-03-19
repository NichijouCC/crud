package sqlx

import (
	"context"
	sqlc "crud/db/sqlc"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Table 通用数据库表操作封装
type Table[T ITable] struct {
	table T             // 表结构实例
	db    *sqlx.DB      // sqlx数据库连接
	Sqlc  *sqlc.Queries // sqlc查询实例
}

// NewModel 创建新的数据库操作模型
func NewModel[T ITable](db *sql.DB) *Table[T] {
	dbx := sqlx.NewDb(db, "mysql")
	var tableInstance T
	return &Table[T]{
		table: tableInstance,
		db:    dbx,
		Sqlc:  sqlc.New(db),
	}
}

// NewModelWithGlobal 使用共享的sqlc/sqlx实例创建新的数据库操作模型
func NewModelWithGlobal[T ITable](db *sql.DB) *Table[T] {
	var tableInstance T
	if _sqlc == nil {
		InitSqlc(_db.DB)
	}
	if _db == nil {
		InitSqlx(db)
	}
	return &Table[T]{
		table: tableInstance,
		db:    _db,
		Sqlc:  _sqlc,
	}
}

// FindAll 查询所有记录
func (m *Table[T]) FindAll(ctx context.Context) ([]T, error) {
	rows, err := FindAll_mysql(ctx, m.db, m.table)
	if err != nil {
		return nil, err
	}
	result := make([]T, len(rows))
	for i, row := range rows {
		result[i] = row.(T)
	}
	return result, nil
}

// FindOneById 根据ID查询单条记录
func (m *Table[T]) FindOneById(ctx context.Context, id int64) (T, error) {
	row, err := FindOneById_mysql(ctx, m.db, m.table, id)
	if err != nil {
		var zero T
		return zero, err
	}
	return row.(T), nil
}

// FindSomeByIds 根据ID列表查询多条记录
func (m *Table[T]) FindSomeByIds(ctx context.Context, ids []int64) ([]T, error) {
	rows, err := FindSomeByIds_mysql(ctx, m.db, m.table, ids)
	if err != nil {
		return nil, err
	}
	result := make([]T, len(rows))
	for i, row := range rows {
		result[i] = row.(T)
	}
	return result, nil
}

// FindSomeByFilter 根据过滤条件查询记录
func (m *Table[T]) FindSomeByFilter(ctx context.Context, filter *QueryFilter) ([]T, error) {
	if filter == nil {
		return m.FindAll(ctx)
	}
	rows, err := FindSomeByFilter_mysql(ctx, m.db, m.table, filter)
	if err != nil {
		return nil, err
	}
	result := make([]T, len(rows))
	for i, row := range rows {
		result[i] = row.(T)
	}
	return result, nil
}

// FindOneByFilter 根据过滤条件查询单条记录
func (m *Table[T]) FindOneByFilter(ctx context.Context, filter *QueryFilter) (T, error) {
	row, err := FindOneByFilter_mysql(ctx, m.db, m.table, filter)
	if err != nil {
		var zero T
		return zero, err
	}
	return row.(T), nil
}

// CreateOne 创建新记录
func (m *Table[T]) CreateOne(ctx context.Context, table ITable) error {
	return CreateOne_mysql(ctx, m.db, table)
}

// UpdateOne 更新记录
func (m *Table[T]) UpdateOne(ctx context.Context, table ITableUpdate) error {
	return UpdateOne_mysql(ctx, m.db, table)
}

// UpdateSomeByIds 更新单条记录
func (m *Table[T]) UpdateSomeByIds(ctx context.Context, table ITableUpdate, ids []int64) error {
	return UpdateSomeByIds_mysql(ctx, m.db, table, ids)
}

// UpdateSomeByFilter 使用过滤条件更新记录
func (m *Table[T]) UpdateSomeByFilter(ctx context.Context, table ITableUpdate, filter *QueryFilter) error {
	return UpdateSomeByFilter_mysql(ctx, m.db, m.table, table, filter)
}

// DeleteOne 删除单条记录
func (m *Table[T]) DeleteOne(ctx context.Context, id int64) error {
	return DeleteOneById_mysql(ctx, m.db, m.table, id)
}

// DeleteSomeByIds 批量删除记录
func (m *Table[T]) DeleteSomeByIds(ctx context.Context, ids []int64) error {
	return DeleteSomeByIds_mysql(ctx, m.db, m.table, ids)
}

// DeleteSomeByFilter 根据过滤条件删除记录
func (m *Table[T]) DeleteSomeByFilter(ctx context.Context, filter *QueryFilter) error {
	return DeleteSomeByFilter_mysql(ctx, m.db, m.table, filter)
}
