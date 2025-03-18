package sqlx

import (
	"errors"
	"strings"
)

var (
	ErrInvalidField    = errors.New("invalid field")
	ErrInvalidOperator = errors.New("invalid operator")
	ErrInvalidLikeChar = errors.New("invalid characters in LIKE value")
	ErrInvalidLikeType = errors.New("LIKE operator only supports string values")
)

// CreateQuerySqlWithFilter 创建查询SQL语句
func CreateQuerySqlWithFilter(table ITable, filter *Filter) (string, []interface{}, error) {
	query := buildBaseSelect(table.TableName())
	if filter == nil {
		return query, nil, nil
	}
	return combineConditions(table, query, nil, filter)
}

// CreateUpdateSqlWithFilter 创建更新SQL语句
func CreateUpdateSqlWithFilter(table ITable, tableUpdate ITableUpdate, filter *Filter) (string, []interface{}, error) {
	query, args, err := buildBaseUpdate(tableUpdate)
	if err != nil {
		return "", nil, err
	}
	if filter == nil {
		return query, nil, nil
	}
	return combineConditions(table, query, args, filter)
}

// CreateDeleteSqlWithFilter 创建删除SQL语句
func CreateDeleteSqlWithFilter(table ITable, filter *Filter) (string, []interface{}, error) {
	query := buildBaseDelete(table.TableName())
	if filter == nil {
		return query, nil, nil
	}
	return combineConditions(table, query, nil, filter)
}

func combineConditions(table ITable, query string, args []interface{}, filter *Filter) (string, []interface{}, error) {
	var builder strings.Builder
	builder.WriteString(query)

	// 处理查询条件
	for i, condition := range filter.Conditions {
		if condition == nil {
			continue
		}
		// 防止SQL注入,验证字段名是否在白名单中
		if _, ok := table.GetAllowedFieldsForFilter()[condition.Field]; !ok {
			return "", nil, ErrInvalidField
		}
		// 添加字段名长度限制
		if len(condition.Field) > 64 {
			return "", nil, errors.New("field name too long")
		}
		// 验证操作符
		upperOperator := strings.ToUpper(condition.Operator)
		if !allowedOperators[upperOperator] {
			return "", nil, ErrInvalidOperator
		}
		condition.Operator = upperOperator

		// 验证 LIKE 操作符的值
		if condition.Operator == "LIKE" {
			if strValue, ok := condition.Value.(string); ok {
				// 更严格的字符检查
				if strings.ContainsAny(strValue, "%_\\'\"`;") || len(strValue) > 100 {
					return "", nil, ErrInvalidLikeChar
				}
				// 添加通配符限制
				if strings.Count(strValue, "%") > 2 {
					return "", nil, errors.New("too many wildcards in LIKE pattern")
				}
				condition.Value = strValue
			} else {
				return "", nil, ErrInvalidLikeType
			}
		}

		if i == 0 {
			builder.WriteString(" WHERE ")
		} else {
			builder.WriteString(" AND ")
		}

		// 使用引号包裹字段名,防止SQL注入
		builder.WriteString("`")
		builder.WriteString(condition.Field)
		builder.WriteString("`")
		builder.WriteString(" ")
		builder.WriteString(condition.Operator)
		builder.WriteString(" ?")
		args = append(args, condition.Value)
	}

	// 验证排序参数
	if filter.Sort != nil {
		// 防止SQL注入,验证字段名是否在白名单中
		if _, ok := table.GetAllowedFieldsForFilter()[filter.Sort.Field]; !ok {
			return "", nil, errors.New("invalid sort field")
		}
		// 验证排序方向
		upperOrder := strings.ToUpper(filter.Sort.Order)
		if upperOrder != "ASC" && upperOrder != "DESC" {
			return "", nil, errors.New("invalid sort order")
		}
		filter.Sort.Order = upperOrder

		builder.WriteString(" ORDER BY `")
		builder.WriteString(filter.Sort.Field)
		builder.WriteString("` ")
		builder.WriteString(filter.Sort.Order)
	}

	// 验证分页参数
	if filter.Pagination != nil {
		if filter.Pagination.Page <= 0 {
			return "", nil, errors.New("页码必须大于0")
		}
		if filter.Pagination.PageSize <= 0 || filter.Pagination.PageSize > 100 {
			return "", nil, errors.New("每页数量必须在1-100之间")
		}

		offset := (filter.Pagination.Page - 1) * filter.Pagination.PageSize
		builder.WriteString(" LIMIT ? OFFSET ?")
		args = append(args, filter.Pagination.PageSize, offset)
	}

	return builder.String(), args, nil
}
