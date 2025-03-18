package middleware

import (
	"crud/pkg/errors"
	"crud/pkg/response"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// ErrorHandler 错误处理中间件
func ErrorHandler() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				HandleError(err, c)
			}
			return err
		}
	}
}

// HandleError 处理错误
func HandleError(err error, c echo.Context) {
	var (
		code     int
		message  string
		status   int
		codeDesc string
	)

	switch e := err.(type) {
	case *errors.Error:
		status = e.HTTPStatus()
		code = int(e.Code)
		message = e.Error()
		codeDesc = errors.GetMessage(e.Code)
	case *echo.HTTPError:
		status = e.Code
		code = int(errors.ErrSystem)
		message = e.Error()
		codeDesc = errors.GetMessage(errors.ErrSystem)
	default:
		status = http.StatusInternalServerError
		code = int(errors.ErrSystem)
		message = "Internal Server Error"
		codeDesc = errors.GetMessage(errors.ErrSystem)
	}

	c.JSON(status, response.Response{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().Unix(),
		CodeDesc:  codeDesc,
	})
}
