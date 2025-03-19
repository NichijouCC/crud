package sqlx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

// buildBaseSelect 构建基本的SELECT查询语句
func buildBaseSelect(tableName string) string {
	return fmt.Sprintf("SELECT * FROM `%s`", tableName)
}

// buildBaseUpdate 构建基本的UPDATE查询语句
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
			placeholders = append(placeholders, fmt.Sprintf("`%s` = ?", fieldName))
			args = append(args, field.Interface())
		}
	}
	if len(placeholders) == 0 {
		return "", nil, errors.New("没有可更新的列")
	}
	query := fmt.Sprintf("UPDATE `%s` SET %s", table.TableName(), strings.Join(placeholders, ", "))
	return query, args, nil
}

// buildBaseDelete 构建基本的DELETE查询语句
func buildBaseDelete(tableName string) string {
	return fmt.Sprintf("DELETE FROM `%s`", tableName)
}

// FindAll_mysql 查询表中的所有记录
func FindAll_mysql(ctx context.Context, db *sqlx.DB, table ITable) ([]ITable, error) {
	query := fmt.Sprintf("SELECT * FROM `%s`", table.TableName())
	var rows []ITable
	if err := db.SelectContext(ctx, &rows, query); err != nil {
		log.Printf("failed to select rows, sql: %s, error: %v", query, err)
		return nil, fmt.Errorf("failed to select rows: %w", err)
	}
	return rows, nil
}

// FindOneById_mysql 根据ID查询单条记录
func FindOneById_mysql(ctx context.Context, db *sqlx.DB, table ITable, id int64) (ITable, error) {
	query := fmt.Sprintf("SELECT * FROM `%s` WHERE `id` = ?", table.TableName())
	var record ITable
	if err := db.GetContext(ctx, &record, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("failed to get row by id, sql: %s, id: %d, error: %v", query, id, err)
		return nil, fmt.Errorf("failed to get row by id: %w", err)
	}
	return record, nil
}

// FindSomeByIds_mysql 根据ID列表批量查询记录
func FindSomeByIds_mysql(ctx context.Context, db *sqlx.DB, table ITable, ids []int64) ([]ITable, error) {
	if len(ids) == 0 {
		return nil, errors.New("ids is empty")
	}
	query, args, err := sqlx.In(fmt.Sprintf("SELECT * FROM `%s` WHERE `id` IN (?)", table.TableName()), ids)
	if err != nil {
		return nil, fmt.Errorf("failed to build IN query: %w", err)
	}
	query = db.Rebind(query)
	var records []ITable
	if err := db.SelectContext(ctx, &records, query, args...); err != nil {
		log.Printf("failed to select rows by ids, sql: %s, args: %v, error: %v", query, args, err)
		return nil, fmt.Errorf("failed to select rows by ids: %w", err)
	}
	return records, nil
}

// FindSomeByFilter_mysql 使用过滤条件查询多条记录
func FindSomeByFilter_mysql(ctx context.Context, db *sqlx.DB, table ITable, filter *QueryFilter) ([]ITable, error) {
	query, args, err := CreateQuerySqlWithFilter(table, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to create filter query: %w", err)
	}
	var records []ITable
	if err := db.SelectContext(ctx, &records, query, args...); err != nil {
		log.Printf("failed to select rows with filter, sql: %s, args: %v, error: %v", query, args, err)
		return nil, fmt.Errorf("failed to select rows with filter: %w", err)
	}
	return records, nil
}

// FindOneByFilter_mysql 使用过滤条件查询单条记录
func FindOneByFilter_mysql(ctx context.Context, db *sqlx.DB, table ITable, filter *QueryFilter) (ITable, error) {
	if filter == nil {
		return nil, fmt.Errorf("filter is nil")
	}
	if filter.Limit != 1 {
		filter.Limit = 1
	}
	query, args, err := CreateQuerySqlWithFilter(table, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to create filter query: %w", err)
	}
	var record ITable
	if err := db.GetContext(ctx, &record, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("failed to get row with filter, sql: %s, args: %v, error: %v", query, args, err)
		return nil, fmt.Errorf("failed to get row with filter: %w", err)
	}
	return record, nil
}

// CreateOne_mysql 创建新记录
func CreateOne_mysql(ctx context.Context, db *sqlx.DB, table ITable) error {
	columns := table.Columns()
	placeholders := make([]string, len(columns))
	quotedColumns := make([]string, len(columns))
	for i, col := range columns {
		placeholders[i] = ":" + col
		quotedColumns[i] = "`" + col + "`"
	}
	query := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)",
		table.TableName(),
		strings.Join(quotedColumns, ", "),
		strings.Join(placeholders, ", "))
	if _, err := db.NamedExecContext(ctx, query, table); err != nil {
		log.Printf("failed to create row, sql: %s, table: %+v, error: %v", query, table, err)
		return fmt.Errorf("failed to create row: %w", err)
	}
	return nil
}

// UpdateOne_mysql 更新单条记录
func UpdateOne_mysql(ctx context.Context, db *sqlx.DB, table ITableUpdate) error {
	query, args, err := buildBaseUpdate(table)
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}
	query += " WHERE `id` = ? LIMIT 1"
	args = append(args, table.GetId())
	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		log.Printf("failed to execute update, sql: %s, args: %v, error: %v", query, args, err)
		return fmt.Errorf("failed to execute update: %w", err)
	}
	return nil
}

// UpdateSomeByIds_mysql 更新单条记录
func UpdateSomeByIds_mysql(ctx context.Context, db *sqlx.DB, table ITableUpdate, ids []int64) error {
	query, args, err := buildBaseUpdate(table)
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}
	query, args, err = sqlx.In(query+" WHERE `id` IN (?)", ids)
	if err != nil {
		return fmt.Errorf("failed to build IN query: %w", err)
	}
	query = db.Rebind(query)
	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		log.Printf("failed to execute update, sql: %s, args: %v, error: %v", query, args, err)
		return fmt.Errorf("failed to execute update: %w", err)
	}
	return nil
}

// UpdateSomeByFilter_mysql 使用过滤条件更新记录
func UpdateSomeByFilter_mysql(ctx context.Context, db *sqlx.DB, table ITable, tableUpdate ITableUpdate, filter *QueryFilter) error {
	if filter == nil {
		return fmt.Errorf("filter is nil")
	}
	query, args, err := CreateUpdateSqlWithFilter(table, tableUpdate, filter)
	if err != nil {
		return fmt.Errorf("failed to create filter update query: %w", err)
	}
	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		log.Printf("failed to execute update with filter, sql: %s, args: %v, error: %v", query, args, err)
		return fmt.Errorf("failed to execute update with filter: %w", err)
	}
	return nil
}

// DeleteOneById_mysql 删除单条记录
func DeleteOneById_mysql(ctx context.Context, db *sqlx.DB, table ITable, id int64) error {
	query := fmt.Sprintf("DELETE FROM `%s` WHERE `id` = ? LIMIT 1", table.TableName())
	if _, err := db.ExecContext(ctx, query, id); err != nil {
		log.Printf("failed to delete row by id, sql: %s, id: %d, error: %v", query, id, err)
		return fmt.Errorf("failed to delete row by id: %w", err)
	}
	return nil
}

// DeleteSomeByIds_mysql 根据ID列表批量删除记录
func DeleteSomeByIds_mysql(ctx context.Context, db *sqlx.DB, table ITable, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("ids is empty")
	}
	query, args, err := sqlx.In(fmt.Sprintf("DELETE FROM `%s` WHERE `id` IN (?)", table.TableName()), ids)
	if err != nil {
		return fmt.Errorf("failed to build IN query: %w", err)
	}
	query = db.Rebind(query)
	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		log.Printf("failed to delete rows by ids, sql: %s, args: %v, error: %v", query, args, err)
		return fmt.Errorf("failed to delete rows by ids: %w", err)
	}
	return nil
}

// DeleteSomeByFilter_mysql 使用过滤条件删除记录
func DeleteSomeByFilter_mysql(ctx context.Context, db *sqlx.DB, table ITable, filter *QueryFilter) error {
	query, args, err := CreateDeleteSqlWithFilter(table, filter)
	if err != nil {
		return fmt.Errorf("failed to create filter delete query: %w", err)
	}
	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		log.Printf("failed to delete rows with filter, sql: %s, args: %v, error: %v", query, args, err)
		return fmt.Errorf("failed to delete rows with filter: %w", err)
	}
	return nil
}
