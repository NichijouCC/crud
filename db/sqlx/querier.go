package sqlx

import (
	"context"
	sqlc "crud/db/sqlc"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// InitSqlx 初始化数据库连接以进行通用的CRUD操作
// 将 *sql.DB 转换为 *sqlx.DB 以支持更多功能
func InitSqlx(db *sql.DB) *sqlx.DB {
	_db = sqlx.NewDb(db, "mysql")
	return _db
}

var _sqlc *sqlc.Queries

// InitSqlc 创建一个新的sqlc查询实例
// 用于生成类型安全的数据库查询
func InitSqlc(db *sql.DB) *sqlc.Queries {
	_sqlc = sqlc.New(db)
	return _sqlc
}

// 全局数据库连接实例
var _db *sqlx.DB

// FindAll 查询表中的所有记录
func FindAll[T ITable](ctx context.Context) ([]T, error) {
	var table T
	rows, err := FindAll_mysql(ctx, _db, table)
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
func FindOneById[T ITable](ctx context.Context, id int64) (T, error) {
	var table T
	row, err := FindOneById_mysql(ctx, _db, table, id)
	if err != nil {
		return table, err
	}
	return row.(T), nil
}

// FindSomeByIds 根据ID列表批量查询记录
func FindSomeByIds[T ITable](ctx context.Context, ids []int64) ([]T, error) {
	var table T
	rows, err := FindSomeByIds_mysql(ctx, _db, table, ids)
	if err != nil {
		return nil, err
	}
	result := make([]T, len(rows))
	for i, row := range rows {
		result[i] = row.(T)
	}
	return result, nil
}

// FindSomeByFilter 使用过滤条件查询多条记录
func FindSomeByFilter[T ITable](ctx context.Context, filter *QueryFilter) ([]T, error) {
	var table T
	rows, err := FindSomeByFilter_mysql(ctx, _db, table, filter)
	if err != nil {
		return nil, err
	}
	result := make([]T, len(rows))
	for i, row := range rows {
		result[i] = row.(T)
	}
	return result, nil
}

// FindOneByFilter 使用过滤条件查询单条记录
func FindOneByFilter[T ITable](ctx context.Context, filter *QueryFilter) (T, error) {
	var table T
	row, err := FindOneByFilter_mysql(ctx, _db, table, filter)
	if err != nil {
		return table, err
	}
	return row.(T), nil
}

// CreateOne 创建新记录
func CreateOne(ctx context.Context, table ITable) error {
	return CreateOne_mysql(ctx, _db, table)
}

// UpdateOne 更新单条记录
func UpdateOne(ctx context.Context, table ITableUpdate) error {
	return UpdateOne_mysql(ctx, _db, table)
}

// UpdateSomeByIds 更新单条记录
func UpdateSomeByIds(ctx context.Context, table ITableUpdate, ids []int64) error {
	return UpdateSomeByIds_mysql(ctx, _db, table, ids)
}

// UpdateSomeByFilter 使用过滤条件更新记录
func UpdateSomeByFilter[T ITable](ctx context.Context, table ITableUpdate, filter *QueryFilter) error {
	var tableInstance T
	return UpdateSomeByFilter_mysql(ctx, _db, tableInstance, table, filter)
}

// DeleteOneById 根据ID删除单条记录
func DeleteOneById[T ITable](ctx context.Context, id int64) error {
	var table T
	return DeleteOneById_mysql(ctx, _db, table, id)
}

// DeleteSomeByIds 根据ID列表批量删除记录
func DeleteSomeByIds[T ITable](ctx context.Context, ids []int64) error {
	var table T
	return DeleteSomeByIds_mysql(ctx, _db, table, ids)
}

// DeleteSomeByFilter 使用过滤条件删除记录
func DeleteSomeByFilter[T ITable](ctx context.Context, filter *QueryFilter) error {
	if filter == nil {
		return fmt.Errorf("过滤条件为空")
	}
	var table T
	return DeleteSomeByFilter_mysql(ctx, _db, table, filter)
}
