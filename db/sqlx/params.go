package sqlx

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Pagination struct {
	Page     int32 `json:"page" query:"page"`           // 页码，从1开始
	PageSize int32 `json:"page_size" query:"page_size"` // 每页数量
}

// Sort 定义排序规则
type Sort struct {
	Field string `json:"field" query:"field"` // 排序字段名
	Order string `json:"order" query:"order"` // 排序方式(ASC/DESC)
}

// Filter 包含分页、排序和过滤条件
type Filter struct {
	Pagination *Pagination  // 分页信息
	Sort       *Sort        // 排序规则
	Conditions []*Condition // 过滤条件列表
}

// Condition 定义单个过滤条件
type Condition struct {
	Field    string      // 字段名
	Value    interface{} // 字段值
	Operator string      // 操作符(=,>,>=,<,<=,LIKE等)
}

var allowedOperators = map[string]bool{"=": true, ">": true, ">=": true, "<": true, "<=": true, "LIKE": true}

// 定义支持的操作符映射
var operatorMap = map[string]string{
	"gt":   ">",
	"gte":  ">=",
	"lt":   "<",
	"lte":  "<=",
	"like": "LIKE",
}

// ParseConditionFromQuery 从URL查询参数解析过滤条件
// 支持的格式:
// field_gt=value -> field > value
// field_gte=value -> field >= value
// field_lt=value -> field < value
// field_lte=value -> field <= value
// field_like=value -> field LIKE value
// field=value -> field = value
func ParseConditionFromQuery(field string, value string) (*Condition, error) {
	// 检查参数
	if field == "" || value == "" {
		return nil, errors.New("field and value cannot be empty")
	}
	// 限制字段名长度
	if len(field) > 64 {
		return nil, errors.New("field name too long")
	}
	parts := strings.Split(field, "_")
	if len(parts) == 1 {
		return &Condition{
			Field:    field,
			Value:    value,
			Operator: "=",
		}, nil
	}
	if len(parts) != 2 {
		return nil, errors.New("invalid field format")
	}

	fieldName := parts[0]
	op := parts[1]
	operator, ok := operatorMap[op]
	if !ok {
		return nil, fmt.Errorf("unsupported operator: %s", op)
	}

	// LIKE操作符特殊处理
	if operator == "LIKE" {
		if strings.ContainsAny(value, "%_\\'\"`;") || len(value) > 100 {
			return nil, ErrInvalidLikeChar
		}
		if strings.Count(value, "%") > 2 {
			return nil, errors.New("too many wildcards in LIKE pattern")
		}
	}

	return &Condition{
		Field:    fieldName,
		Value:    value,
		Operator: operator,
	}, nil
}

// ParseConditionsFromQuery 从URL查询参数解析出查询条件
func ParseConditionsFromQuery(query map[string][]string) ([]*Condition, error) {
	var conditions []*Condition
	for field, values := range query {
		// 跳过空值
		if len(values) == 0 || values[0] == "" {
			continue
		}
		// 只取第一个值
		value := values[0]
		condition, err := ParseConditionFromQuery(field, value)
		if err != nil {
			return nil, fmt.Errorf("解析条件失败 %s: %v", field, err)
		}
		conditions = append(conditions, condition)
	}
	return conditions, nil
}

// ParseFilterFromContext 从echo.Context解析出Filter
func ParseFilterFromContext(params map[string][]string) (*Filter, error) {
	// 获取所有查询参数
	// 解析查询条件
	conditions, err := ParseConditionsFromQuery(params)
	if err != nil {
		return nil, err
	}

	// 解析分页参数
	var pagination *Pagination
	var page, pageSize int32 = 1, 10 // 默认值
	if pageValues, ok := params["page"]; ok && len(pageValues) > 0 {
		if p, err := strconv.ParseInt(pageValues[0], 10, 32); err == nil && p > 0 {
			page = int32(p)
			pagination = &Pagination{
				Page:     page,
				PageSize: pageSize,
			}
		}
	}
	if sizeValues, ok := params["page_size"]; ok && len(sizeValues) > 0 {
		if s, err := strconv.ParseInt(sizeValues[0], 10, 32); err == nil && s > 0 {
			pageSize = int32(s)
			if pagination == nil {
				pagination = &Pagination{
					Page:     page,
					PageSize: pageSize,
				}
			}
		}
	}
	// 解析排序
	var sort *Sort
	if fieldValues, ok := params["sort_field"]; ok && len(fieldValues) > 0 {
		sort = &Sort{
			Field: fieldValues[0],
			Order: "ASC", // 默认升序
		}
		if orderValues, ok := params["sort_order"]; ok && len(orderValues) > 0 {
			if strings.ToUpper(orderValues[0]) == "DESC" {
				sort.Order = "DESC"
			}
		}
	}

	if len(conditions) != 0 || pagination != nil || sort != nil {
		return &Filter{
			Conditions: conditions,
			Pagination: pagination,
			Sort:       sort,
		}, nil
	}
	return nil, nil
}
