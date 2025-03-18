package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
)

// Recover 恢复中间件
func Recover() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					c.Error(r.(error))
				}
			}()
			return next(c)
		}
	}
}

// CORS 跨域中间件
func CORS() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")
			c.Response().Header().Set(echo.HeaderAccessControlAllowMethods, "GET, POST, PUT, DELETE, OPTIONS")
			c.Response().Header().Set(echo.HeaderAccessControlAllowHeaders, "Content-Type, Authorization")
			return next(c)
		}
	}
}

// Logger 日志中间件
func Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			path := req.URL.Path
			if path == "" {
				path = "/"
			}

			start := time.Now()
			logger := c.Logger()

			// 记录请求开始
			logger.Printf("请求开始: %s %s", req.Method, path)

			// 处理请求
			err := next(c)

			// 记录请求完成
			logger.Printf("请求完成: %s %s %d %s",
				req.Method,
				path,
				c.Response().Status,
				time.Since(start).String(),
			)

			return err
		}
	}
}
