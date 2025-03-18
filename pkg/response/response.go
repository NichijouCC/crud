package response

import (
	"crud/pkg/errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// Response 统一响应结构
type Response struct {
	Code      int         `json:"code"`            // 业务错误码
	CodeDesc  string      `json:"code_desc"`       // 错误码说明
	Message   string      `json:"message"`         // 响应消息
	Data      interface{} `json:"data,omitempty"`  // 响应数据
	Timestamp int64       `json:"timestamp"`       // 时间戳
	Error     string      `json:"error,omitempty"` // 错误信息
}

// Success 成功响应
func Success(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, Response{
		Code:      int(errors.Success),
		CodeDesc:  errors.GetMessage(errors.Success),
		Data:      data,
		Timestamp: time.Now().Unix(),
	})
}

// Error 错误响应
func Error(code errors.ErrorCode, err error, message string) *errors.Error {
	return errors.Wrap(code, message, err)
}

// 系统级错误
func SystemError(err error) *errors.Error {
	return Error(errors.ErrSystem, err, "")
}

func DatabaseError(err error) *errors.Error {
	return Error(errors.ErrDatabase, err, "")
}

func CacheError(err error) *errors.Error {
	return Error(errors.ErrCache, err, "")
}

func NetworkError(err error) *errors.Error {
	return Error(errors.ErrNetwork, err, "")
}

func ServiceUnavailableError(err error) *errors.Error {
	return Error(errors.ErrServiceUnavailable, err, "")
}

// 业务级错误
func BadRequest(err error) *errors.Error {
	return Error(errors.ErrInvalidParams, err, "")
}

func ValidationError(err error) *errors.Error {
	return Error(errors.ErrValidation, err, "")
}

func BusinessError(err error) *errors.Error {
	return Error(errors.ErrBusiness, err, "")
}

func NotFound(err error) *errors.Error {
	return Error(errors.ErrNotFound, err, "")
}

func DuplicateError(err error) *errors.Error {
	return Error(errors.ErrDuplicate, err, "")
}

func ConflictError(err error) *errors.Error {
	return Error(errors.ErrConflict, err, "")
}

// 权限级错误
func Unauthorized(err error) *errors.Error {
	return Error(errors.ErrUnauthorized, err, "")
}

func Forbidden(err error) *errors.Error {
	return Error(errors.ErrForbidden, err, "")
}

func ExpiredError(err error) *errors.Error {
	return Error(errors.ErrExpired, err, "")
}

func InvalidTokenError(err error) *errors.Error {
	return Error(errors.ErrInvalidToken, err, "")
}
