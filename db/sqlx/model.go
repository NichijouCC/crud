package sqlx

import (
	"context"
	sqlc "crud/db/sqlc"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

// Model 通用数据库操作模型
type Model[T ITable] struct {
	table T             // 表结构实例
	db    *sqlx.DB      // sqlx数据库连接
	Sqlc  *sqlc.Queries // sqlc查询实例
}

// NewModel 创建新的数据库操作模型
func NewModel[T ITable](db *sql.DB) *Model[T] {
	dbx := sqlx.NewDb(db, "mysql")
	var tableInstance T
	return &Model[T]{
		table: tableInstance,
		db:    dbx,
		Sqlc:  sqlc.New(db),
	}
}

// NewModelWithCommonSqlC 使用共享的sqlc实例创建新的数据库操作模型
func NewModelWithCommonSqlC[T ITable](db *sql.DB) *Model[T] {
	dbx := sqlx.NewDb(db, "mysql")
	var tableInstance T
	if _sqlc == nil {
		InitSqlcQueries(_db.DB)
	}
	return &Model[T]{
		table: tableInstance,
		db:    dbx,
		Sqlc:  _sqlc,
	}
}

// FindRows 查询所有记录
func (m *Model[T]) FindRows(ctx context.Context) ([]*T, error) {
	tableName := m.table.TableName()
	query := "SELECT * FROM " + tableName
	var rows []*T
	err := m.db.SelectContext(ctx, &rows, query)
	return rows, err
}

// FindRowsByIDs 根据ID列表查询多条记录
func (m *Model[T]) FindRowsByIDs(ctx context.Context, ids []int64) ([]*T, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("ids is empty")
	}
	query := "SELECT * FROM " + m.table.TableName() + " WHERE id IN (?) "
	var records []*T
	err := m.db.SelectContext(ctx, &records, query, ids)
	return records, err
}

// FindRowsByFilter 根据过滤条件查询记录
func (m *Model[T]) FindRowsByFilter(ctx context.Context, filter *Filter) ([]*T, error) {
	if filter == nil {
		return m.FindRows(ctx)
	}
	query, args, err := CreateQuerySqlWithFilter(m.table, filter)
	if err != nil {
		return nil, err
	}
	var records []*T
	err = m.db.SelectContext(ctx, &records, query, args...)
	return records, err
}

// FindByID 根据ID查询单条记录
func (m *Model[T]) FindByID(ctx context.Context, id int64) (*T, error) {
	query := "SELECT * FROM " + m.table.TableName() + " WHERE id = ?"
	err := m.db.SelectContext(ctx, &m.table, query, id)
	return &m.table, err
}

// FindByFilter 根据过滤条件查询单条记录
func (m *Model[T]) FindByFilter(ctx context.Context, filter *Filter) (*T, error) {
	query, args, err := CreateQuerySqlWithFilter(m.table, filter)
	if err != nil {
		return nil, err
	}
	var record *T
	err = m.db.SelectContext(ctx, &record, query, args...)
	return record, err
}

// Create 创建新记录
func (m *Model[T]) Create(ctx context.Context, table ITable) error {
	columns := table.Columns()
	placeholders := make([]string, len(columns))
	for i := range columns {
		placeholders[i] = ":" + columns[i]
	}

	query := "INSERT INTO " + m.table.TableName() + " (" + strings.Join(columns, ", ") + ") VALUES (" + strings.Join(placeholders, ", ") + ")"
	_, err := m.db.NamedExecContext(ctx, query, table)
	return err
}

// Update 更新记录
func (m *Model[T]) Update(ctx context.Context, table ITableUpdate) error {
	columns := table.Columns()
	var placeholders []string
	v := reflect.ValueOf(table).Elem()
	for i := range columns {
		if v.Field(i).IsZero() {
			continue
		}
		placeholders = append(placeholders, columns[i]+" = :"+columns[i])
	}
	query := "UPDATE " + m.table.TableName() + " SET " + strings.Join(placeholders, ", ") + " WHERE id = :id"
	_, err := m.db.NamedExecContext(ctx, query, table)
	return err
}

// Delete 删除单条记录
func (m *Model[T]) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM " + m.table.TableName() + " WHERE id = ?"
	_, err := m.db.ExecContext(ctx, query, id)
	return err
}

// DeleteRowsByIDs 批量删除记录
func (m *Model[T]) DeleteRowsByIDs(ctx context.Context, ids []int64) error {
	query := "DELETE FROM " + m.table.TableName() + " WHERE id IN (?) "
	_, err := m.db.ExecContext(ctx, query, ids)
	return err
}

// DeleteRowsByFilter 根据过滤条件删除记录
func (m *Model[T]) DeleteRowsByFilter(ctx context.Context, filter *Filter) error {
	query, args, err := CreateQuerySqlWithFilter(m.table, filter)
	if err != nil {
		return err
	}
	_, err = m.db.ExecContext(ctx, query, args...)
	return err
}
