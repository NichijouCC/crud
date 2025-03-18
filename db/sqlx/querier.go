package sqlx

import (
	"context"
	sqlc "crud/db/sqlc"
	"database/sql"
	"errors"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

// InitSqlx 初始化通用CRUD操作的数据库连接
// 将 *sql.DB 转换为 *sqlx.DB 以支持更多功能
func InitSqlx(db *sql.DB) *sqlx.DB {
	_db = sqlx.NewDb(db, "mysql")
	return _db
}

var _sqlc *sqlc.Queries

// InitSqlcQueries 创建一个新的sqlc查询实例
// 用于生成类型安全的数据库查询
func InitSqlcQueries(db *sql.DB) *sqlc.Queries {
	_sqlc = sqlc.New(db)
	return _sqlc
}

type ITable interface {
	TableName() string
	GetAllowedFieldsForFilter() map[string]struct{}
	Columns() []string
	GetId() int64
}

type ITableUpdate interface {
	TableName() string
	Columns() []string
	GetId() int64
}

var _db *sqlx.DB

func buildBaseSelect(tableName string) string {
	return "SELECT * FROM " + tableName
}

func buildBaseUpdate(table ITableUpdate) (string, []interface{}, error) {
	v := reflect.ValueOf(table).Elem()
	t := v.Type()

	var placeholders []string
	var args []interface{}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.IsZero() {
			fieldName := t.Field(i).Name
			if fieldName == "Id" {
				continue
			}
			fieldValue := field.Interface()
			placeholders = append(placeholders, fieldName+" = ?"+fieldName)
			args = append(args, fieldValue)
		}
	}
	if len(placeholders) == 0 {
		return "", nil, errors.New("no updateable columns")
	}
	query := "UPDATE " + table.TableName() + " SET " + strings.Join(placeholders, ", ")
	return query, args, nil
}

func buildBaseDelete(tableName string) string {
	return "DELETE FROM " + tableName
}

func FindRows[T ITable](ctx context.Context) ([]*T, error) {
	var tableInstance T
	tableName := tableInstance.TableName()
	query := buildBaseSelect(tableName)
	var rows []*T
	err := _db.SelectContext(ctx, &rows, query)
	return rows, err
}

func FindRowsWithFilter[T ITable](ctx context.Context, filter *Filter) ([]*T, error) {
	var tableInstance T
	query, args, err := CreateQuerySqlWithFilter(tableInstance, filter)
	if err != nil {
		return nil, err
	}
	var rows []*T
	err = _db.SelectContext(ctx, &rows, query, args...)
	return rows, err
}

func Find[T ITable](ctx context.Context) (*T, error) {
	var tableInstance T
	query := "SELECT * FROM " + tableInstance.TableName() + " WHERE id = ?"
	err := _db.SelectContext(ctx, &tableInstance, query, tableInstance.GetId())
	return &tableInstance, err
}

func FindRowsByIds[T ITable](ctx context.Context, ids []int64) ([]*T, error) {
	var tableInstance T
	var records []*T
	query := "SELECT * FROM " + tableInstance.TableName() + " WHERE id IN (?) "
	err := _db.SelectContext(ctx, &records, query, ids)
	return records, err
}

func FindRowsByFilter[T ITable](ctx context.Context, filter *Filter) ([]*T, error) {
	var tableInstance T
	var records []*T
	query, args, err := CreateQuerySqlWithFilter(tableInstance, filter)
	if err != nil {
		return nil, err
	}
	err = _db.SelectContext(ctx, &records, query, args...)
	return records, err
}

func FindByFilter[T ITable](ctx context.Context, filter *Filter) (*T, error) {
	var tableInstance T
	query, args, err := CreateQuerySqlWithFilter(tableInstance, filter)
	if err != nil {
		return &tableInstance, err
	}
	err = _db.SelectContext(ctx, &tableInstance, query, args...)
	return &tableInstance, err
}

func Create(ctx context.Context, table ITable) error {
	columns := table.Columns()
	placeholders := make([]string, len(columns))
	for i := range columns {
		placeholders[i] = ":" + columns[i]
	}
	query := "INSERT INTO " + table.TableName() + " (" + strings.Join(columns, ", ") + ") VALUES (" + strings.Join(placeholders, ", ") + ")"
	_, err := _db.NamedExecContext(ctx, query, table)
	return err
}

func Update(ctx context.Context, table ITableUpdate) error {
	query, args, err := buildBaseUpdate(table)
	if err != nil {
		return err
	}
	query += " WHERE id = ?"
	args = append(args, table.GetId())
	result := _db.MustExecContext(ctx, query, args...)
	_, err = result.RowsAffected()
	return err
}

func UpdateWithFilter[T ITable](ctx context.Context, table ITableUpdate, filter *Filter) error {
	var tableInstance T
	query, args, err := CreateUpdateSqlWithFilter(tableInstance, table, filter)
	if err != nil {
		return err
	}
	result := _db.MustExecContext(ctx, query, args...)
	_, err = result.RowsAffected()
	return err
}

func Delete[T ITable](ctx context.Context, id int64) error {
	var tableInstance T
	query := "DELETE FROM " + tableInstance.TableName() + " WHERE id = ?"
	_, err := _db.ExecContext(ctx, query, id)
	return err
}

func DeleteRowsByIds[T ITable](ctx context.Context, ids []int64) error {
	var tableInstance T
	query := "DELETE FROM " + tableInstance.TableName() + " WHERE id IN (?) "
	_, err := _db.ExecContext(ctx, query, ids)
	return err
}

func DeleteRowsByFilter[T ITable](ctx context.Context, filter *Filter) error {
	var tableInstance T
	query, args, err := CreateDeleteSqlWithFilter(tableInstance, filter)
	if err != nil {
		return err
	}
	_, err = _db.ExecContext(ctx, query, args...)
	return err
}
