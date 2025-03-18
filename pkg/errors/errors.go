package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode 错误码类型
type ErrorCode int

// 成功状态
const Success ErrorCode = 0

// 系统级错误码 (1000-1999)
const (
	ErrSystem ErrorCode = iota + 1000
	ErrDatabase
	ErrCache
	ErrNetwork
	ErrServiceUnavailable
)

// 业务级错误码 (2000-2999)
const (
	ErrInvalidParams ErrorCode = iota + 2000
	ErrValidation
	ErrBusiness
	ErrNotFound
	ErrDuplicate
	ErrConflict
)

// 权限级错误码 (3000-3999)
const (
	ErrUnauthorized ErrorCode = iota + 3000
	ErrForbidden
	ErrExpired
	ErrInvalidToken
)

// messages 错误码对应的中文消息映射
var messages = map[ErrorCode]string{
	Success:               "成功",
	ErrSystem:             "系统错误",
	ErrDatabase:           "数据库错误",
	ErrCache:              "缓存错误",
	ErrNetwork:            "网络错误",
	ErrServiceUnavailable: "服务不可用",
	ErrInvalidParams:      "无效的参数",
	ErrValidation:         "验证错误",
	ErrBusiness:           "业务错误",
	ErrNotFound:           "未找到",
	ErrDuplicate:          "重复",
	ErrConflict:           "冲突",
	ErrUnauthorized:       "未授权",
	ErrForbidden:          "禁止访问",
	ErrExpired:            "已过期",
	ErrInvalidToken:       "无效的令牌",
}

// GetMessage 获取错误码对应的中文消息
func GetMessage(code ErrorCode) string {
	if msg, ok := messages[code]; ok {
		return msg
	}
	return "未知错误"
}

// Error 统一错误结构
type Error struct {
	Code    ErrorCode // 错误码
	Message string    // 错误消息
	Err     error     // 原始错误
}

// New 创建新的错误
func New(code ErrorCode, message string) *Error {
	if message == "" {
		message = GetMessage(code)
	}
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Wrap 包装已有错误
func Wrap(code ErrorCode, message string, err error) *Error {
	if message == "" {
		message = GetMessage(code)
	}
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Error 实现 error 接口
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap 实现 errors.Unwrap 接口
func (e *Error) Unwrap() error {
	return e.Err
}

// GetMessage 获取错误消息
func (e *Error) GetMessage() string {
	return e.Message
}

// HTTPStatus 获取 HTTP 状态码
func (e *Error) HTTPStatus() int {
	switch e.Code {
	case Success:
		return http.StatusOK
	case ErrSystem, ErrDatabase:
		return http.StatusInternalServerError
	case ErrCache:
		return http.StatusServiceUnavailable
	case ErrNetwork:
		return http.StatusBadGateway
	case ErrServiceUnavailable:
		return http.StatusServiceUnavailable
	case ErrInvalidParams:
		return http.StatusBadRequest
	case ErrValidation:
		return http.StatusUnprocessableEntity
	case ErrBusiness:
		return http.StatusBadRequest
	case ErrNotFound:
		return http.StatusNotFound
	case ErrDuplicate:
		return http.StatusConflict
	case ErrConflict:
		return http.StatusConflict
	case ErrUnauthorized, ErrExpired, ErrInvalidToken:
		return http.StatusUnauthorized
	case ErrForbidden:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

// IsNotFound 判断是否是未找到错误
func IsNotFound(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == ErrNotFound
	}
	return false
}

// IsValidationError 判断是否是验证错误
func IsValidationError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == ErrValidation
	}
	return false
}

// IsBusinessError 判断是否是业务错误
func IsBusinessError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == ErrBusiness
	}
	return false
}

// IsUnauthorized 判断是否是未授权错误
func IsUnauthorized(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == ErrUnauthorized
	}
	return false
}

// IsForbidden 判断是否是禁止访问错误
func IsForbidden(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == ErrForbidden
	}
	return false
}
