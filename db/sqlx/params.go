package sqlx

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// QueryFilter 包含分页、排序和过滤条件
type QueryFilter struct {
	Conditions []*QueryCondition // 过滤条件列表
	Limit      int               // 限制返回的记录数
	Offset     int               // 偏移量
	SortField  string            // 排序字段名
	SortOrder  string            // 排序方式(ASC/DESC)
}

// QueryCondition 定义单个过滤条件
type QueryCondition struct {
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

// ParseQueryConditionFromUrlParam 从URL查询参数解析过滤条件
// 支持的格式:
// field_gt=value -> field > value
// field_gte=value -> field >= value
// field_lt=value -> field < value
// field_lte=value -> field <= value
// field_like=value -> field LIKE value
// field=value -> field = value
func ParseQueryConditionFromUrlParam(field string, value string) (*QueryCondition, error) {
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
		return &QueryCondition{
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

	return &QueryCondition{
		Field:    fieldName,
		Value:    value,
		Operator: operator,
	}, nil
}

// ParseQueryConditionsFromUrlParams 从URL查询参数解析出查询条件
func ParseQueryConditionsFromUrlParams(query map[string][]string) ([]*QueryCondition, error) {
	var conditions []*QueryCondition
	for field, values := range query {
		// 跳过空值
		if len(values) == 0 || values[0] == "" {
			continue
		}
		// 只取第一个值
		value := values[0]
		condition, err := ParseQueryConditionFromUrlParam(field, value)
		if err != nil {
			return nil, fmt.Errorf("解析条件失败 %s: %v", field, err)
		}
		conditions = append(conditions, condition)
	}
	return conditions, nil
}

// ParseQueryFilterFromUrlParams 从echo.Context解析出Filter
func ParseQueryFilterFromUrlParams(params map[string][]string) (*QueryFilter, error) {
	conditions, err := ParseQueryConditionsFromUrlParams(params)
	if err != nil {
		return nil, err
	}

	// 解析分页参数
	var limit, offset int = 0, 0

	if pageValues, ok := params["page"]; ok && len(pageValues) > 0 {
		if p, err := strconv.ParseInt(pageValues[0], 10, 32); err == nil && p > 0 {
			page := int32(p)
			var pageSize int32 = 10 // 默认值
			if sizeValues, ok := params["page_size"]; ok && len(sizeValues) > 0 {
				if s, err := strconv.ParseInt(sizeValues[0], 10, 32); err == nil && s > 0 {
					pageSize = int32(s)
				}
			}
			limit = int(pageSize)
			offset = int((page - 1) * pageSize)
		}
	}

	// 解析排序
	var sortField, sortOrder string = "", ""
	if fieldValues, ok := params["sort_field"]; ok && len(fieldValues) > 0 {
		sortField = fieldValues[0]
		sortOrder = "ASC"
		if orderValues, ok := params["sort_order"]; ok && len(orderValues) > 0 {
			if strings.ToUpper(orderValues[0]) == "DESC" {
				sortOrder = "DESC"
			}
		}
	}

	if len(conditions) != 0 || limit != 0 || offset != 0 || sortField != "" || sortOrder != "" {
		return &QueryFilter{
			Conditions: conditions,
			Limit:      limit,
			Offset:     offset,
			SortField:  sortField,
			SortOrder:  sortOrder,
		}, nil
	}
	return nil, nil
}

// 在查询参数中标记需要的字段
const RequiredFieldsKey = "atts_require"

// 在查询参数中标记省略的字段
const OmittedFieldsKey = "atts_omit"

// FieldFilter 结构体用于存储需要和省略的字段
type FieldFilter struct {
	RequiredFields []string // 需要的字段列表
	OmittedFields  []string // 省略的字段列表
}

// ParseFieldFilterFromQuery 从查询参数中解析出 FieldFilter
func ParseFieldFilterFromQuery(params map[string][]string) (*FieldFilter, bool) {
	for key, values := range params {
		if key == RequiredFieldsKey {
			fieldFilter := &FieldFilter{}
			// 解析需要的字段
			for _, value := range values {
				fieldFilter.RequiredFields = append(fieldFilter.RequiredFields, strings.Split(value, ",")...)
			}
			return fieldFilter, true
		}
		if key == OmittedFieldsKey {
			fieldFilter := &FieldFilter{}
			// 解析省略的字段
			for _, value := range values {
				fieldFilter.OmittedFields = append(fieldFilter.OmittedFields, strings.Split(value, ",")...)
			}
			return fieldFilter, true
		}
	}
	return nil, false
}

func BuildSelectWithFieldFilter(table ITable, fieldFilter *FieldFilter) (string, []interface{}, error) {
	allowedFields := table.ColumnsMap()
	var selectedFields []string

	// 过滤需要的字段
	if fieldFilter != nil {
		if len(fieldFilter.RequiredFields) > 0 {
			for _, field := range fieldFilter.RequiredFields {
				if _, ok := allowedFields[field]; ok {
					selectedFields = append(selectedFields, field)
				}
			}
		} else {
			for field := range allowedFields {
				selectedFields = append(selectedFields, field)
			}
		}

		// 过滤省略的字段
		if len(fieldFilter.OmittedFields) > 0 {
			var filteredFields []string
			omitMap := make(map[string]struct{})
			for _, field := range fieldFilter.OmittedFields {
				omitMap[field] = struct{}{}
			}
			for _, field := range selectedFields {
				if _, ok := omitMap[field]; !ok {
					filteredFields = append(filteredFields, field)
				}
			}
			selectedFields = filteredFields
		}
	} else {
		for field := range allowedFields {
			selectedFields = append(selectedFields, field)
		}
	}

	if len(selectedFields) == 0 {
		return "", nil, fmt.Errorf("no valid fields selected")
	}

	query := fmt.Sprintf("SELECT %s FROM `%s`", strings.Join(selectedFields, ", "), table.TableName())
	return query, nil, nil
}
